package qstore

import (
	"gat/internal/config"
	"os"
	"path/filepath"
	"time"

	"gat/internal/qstore"
)

func pingQdrant(cfg config.AppConfig) error {
	st, err := qstore.New(cfg.QdrantAddr, cfg.CollectionName, cfg.EmbeddingDim)
	if err != nil {
		return err
	}
	defer st.Close()

	ctx, cancel := qstore.CtxTimeout(2 * time.Second)
	defer cancel()
	return st.Health(ctx)
}

// ensureCollection makes sure the target collection exists with the expected vector size.
func ensureCollection(cfg config.AppConfig) error {
	st, err := qstore.New(cfg.QdrantAddr, cfg.CollectionName, cfg.EmbeddingDim)
	if err != nil {
		return err
	}
	defer st.Close()

	ctx, cancel := qstore.CtxTimeout(5 * time.Second)
	defer cancel()
	return st.EnsureCollection(ctx)
}

func defaultQdrantDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".llm-pyhelp", "qdrant")
}
