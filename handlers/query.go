package handlers

import (
	"net/http"

	"agentcore/config"
	"agentcore/db"
	"agentcore/models"
	"agentcore/repository"
	routerservice "agentcore/services/router"
	"github.com/gin-gonic/gin"
)

var (
	routerServiceFactory = func() routeEngine {
		return routerservice.NewService(
			config.RouterEnableLLM,
			config.RouterRulesConfidenceThreshold,
			nil,
		)
	}

	routingDecisionRepositoryFactory = func() repository.RoutingDecisionRepository {
		if db.DB == nil {
			return nil
		}
		return repository.NewRoutingDecisionRepository(db.DB)
	}
)

type routeEngine interface {
	Route(req models.QueryRequest) (models.RouterDecision, error)
}

var queryResponder = func(req models.QueryRequest, decision models.RouterDecision) (models.QueryResponse, error) {
	return models.QueryResponse{
		ConversationID: req.Context.ConversationID,
		Answer:         "Query received: " + req.Prompt,
		Intent:         decision.Intent,
		Confidence:     decision.Confidence,
		Reason:         decision.Reason,
	}, nil
}

func Query(c *gin.Context) {
	var req models.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteValidationError(c, "Invalid query request", err.Error())
		return
	}

	decision, err := routerServiceFactory().Route(req)
	if err != nil {
		WriteInternalError(c, "Failed to classify query")
		return
	}

	if repo := routingDecisionRepositoryFactory(); repo != nil {
		_, err := repo.Create(&models.RoutingDecisionRecord{
			ConversationID: req.Context.ConversationID,
			Intent:         decision.Intent,
			Confidence:     decision.Confidence,
			Reason:         decision.Reason,
			Classifier:     decision.Classifier,
			PolicyVersion:  decision.PolicyVersion,
			TraceID:        req.Context.TraceID,
			UserID:         req.Context.UserID,
			Source:         req.Context.Source,
			Tags:           req.Context.Tags,
		})
		if err != nil {
			WriteInternalError(c, "Failed to persist routing decision")
			return
		}
	}

	resp, err := queryResponder(req, decision)
	if err != nil {
		WriteInternalError(c, "Failed to process query")
		return
	}

	c.JSON(http.StatusOK, resp)
}
