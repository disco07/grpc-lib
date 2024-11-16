package files

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/metadata"
)

type Metadata struct {
	Name  string    `form:"name" json:"name"`
	Age   int       `form:"age" json:"age"`
	Price float64   `form:"price" json:"price"`
	Date  time.Time `form:"date" json:"date"`
}

// Exemple de struct pour mapper les données du formulaire
type TestStruct struct {
	Metadata Metadata                `form:"metadata"`
	Files    []*multipart.FileHeader `form:"files"`
}

func TestParseMultipartForm(t *testing.T) {
	// Création d'un corps multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	metadataWriter, err := writer.CreateFormField("metadata")
	assert.NoError(t, err)

	data := `{"name": "John Doe", "age":30, "price":100.0, "date":"2000-01-01T00:00:00Z"}`
	_, err = metadataWriter.Write([]byte(data))
	assert.NoError(t, err)

	// Ajouter un fichier
	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", `form-data; name="files"; filename="test.txt"`)
	partHeader.Set("Content-Type", "text/plain")
	part, err := writer.CreatePart(partHeader)
	assert.NoError(t, err)
	part.Write([]byte("File content"))

	// Fermer le writer
	assert.NoError(t, writer.Close())

	// Créer un Content-Type avec boundary
	contentType := writer.FormDataContentType()

	// Injecter Content-Type dans le contexte gRPC
	md := metadata.New(map[string]string{
		fmt.Sprintf("%scontent-type", runtime.MetadataPrefix): contentType,
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Extraire la boundary depuis le contexte
	boundary, err := extractBoundaryFromContext(ctx)
	assert.NoError(t, err, "L'extraction de la boundary devrait réussir")
	assert.NotEmpty(t, boundary, "La boundary ne doit pas être vide")

	// Créer l'objet HttpBody
	httpBody := &httpbody.HttpBody{
		ContentType: contentType,
		Data:        body.Bytes(),
	}

	// Appeler ParseMultipartForm
	var result *TestStruct
	result, err = ParseMultipartForm[TestStruct](ctx, httpBody)

	// Vérifier que l'erreur est nulle
	if err != nil {
		t.Fatalf("ParseMultipartForm a échoué : %v", err)
	}

	// Vérifications des résultats
	assert.Equal(t, "John Doe", result.Metadata.Name, "Le champ Name doit être mappé correctement")
	assert.Equal(t, 30, result.Metadata.Age, "Le champ Age doit être mappé correctement")
	assert.Equal(t, 100.0, result.Metadata.Price, "Le champ Price doit être mappé correctement")
	assert.Equal(t, "2000-01-01T00:00:00Z", result.Metadata.Date.Format(time.RFC3339), "Le champ Date doit être mappé correctement")
	if assert.NotNil(t, result.Files) && assert.Len(t, result.Files, 1) {
		assert.Equal(t, "test.txt", result.Files[0].Filename, "Le fichier doit être correctement identifié")
	}
}
