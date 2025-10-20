package qstore

import (
	"context"

	qdrant "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// dado relacionado ao vetor
type Item struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type Store struct {
	points     qdrant.PointsClient
	colls      qdrant.CollectionsClient
	collection string
	model      string
	dim        int
	conn       *grpc.ClientConn
}

func New(address, collection string, model string, dim int) (*Store, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Store{
		points:     qdrant.NewPointsClient(conn),
		colls:      qdrant.NewCollectionsClient(conn),
		collection: collection,
		model:      model,
		dim:        dim,
		conn:       conn,
	}, nil
}

// encerra conexao com o qdrant
func (s *Store) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// cria collection se nao existir
func (s *Store) EnsureCollection(ctx context.Context) error {
	if _, err := s.colls.Get(ctx, &qdrant.GetCollection{
		CollectionName: s.collection,
	}); err == nil {
		return nil
	}

	_, err := s.colls.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: s.collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     int32(s.dim),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	return err
}
