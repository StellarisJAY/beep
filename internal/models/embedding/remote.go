package embedding

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
)

// remoteEmbedder openai-compatible api
type remoteEmbedder struct {
	client     openai.Client
	modelName  string
	dimensions int64
}

func (r *remoteEmbedder) EmbedString(ctx context.Context, text string) ([]float64, error) {
	response, err := r.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input:          openai.EmbeddingNewParamsInputUnion{OfString: param.NewOpt(text)},
		Model:          r.modelName,
		Dimensions:     param.NewOpt(r.dimensions),
		EncodingFormat: "float",
	})
	if err != nil {
		return nil, err
	}
	return response.Data[0].Embedding, nil
}

func (r *remoteEmbedder) BatchEmbedStrings(ctx context.Context, texts []string) ([][]float64, error) {
	response, err := r.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input:          openai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: texts},
		Model:          r.modelName,
		Dimensions:     param.NewOpt(r.dimensions),
		EncodingFormat: "float",
	})
	if err != nil {
		return nil, err
	}
	embeddings := make([][]float64, 0, len(response.Data))
	for _, data := range response.Data {
		embeddings = append(embeddings, data.Embedding)
	}
	return embeddings, nil
}
