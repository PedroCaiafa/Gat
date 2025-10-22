package qstore

import (
	"context"
	"time"

	"github.com/qdrant/go-client/qdrant"

	"gat/internal/config"
)

// dado relacionado ao vetor
type Item struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type QStore struct {
	client     *qdrant.Client
	collection string
}

func CtxTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}

func NewQStore(cfg config.AppConfig) (*QStore, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.QdrantAddress,
		Port:   cfg.QdrantPort,
		APIKey: cfg.QdrantAPIKey,
	})
	if err != nil {
		return nil, err
	}
	return &QStore{
		client:     client,
		collection: cfg.CollectionName,
	}, nil
}

func ensureCollection(cfg config.AppConfig) error {
	st, err := NewQStore(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := CtxTimeout(10 * time.Second)
	defer cancel()
	exists, err := st.client.CollectionExists(ctx, cfg.CollectionName)
	if err != nil {
		return err
	}
	if !exists {
		return st.client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: cfg.CollectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size: cfg.EmbeddingDim,
			}),
		})
	}
	return nil
}

/*func defaultQdrantDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".llm-pyhelp", "qdrant")
}
*/
