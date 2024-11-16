package files

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/metadata"
)

// Exemple de struct pour mapper les données du formulaire
type TestStruct struct {
	Name  string                  `form:"name"`
	Age   int                     `form:"age"`
	Files []*multipart.FileHeader `form:"files"`
}

func TestParseMultipartForm(t *testing.T) {
	// Création d'un corps multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Ajouter un champ texte
	err := writer.WriteField("name", "John Doe")
	assert.NoError(t, err)

	// Ajouter un champ numérique
	err = writer.WriteField("age", "30")
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
	assert.Equal(t, "John Doe", result.Name, "Le champ Name doit être mappé correctement")
	assert.Equal(t, 30, result.Age, "Le champ Age doit être mappé correctement")
	if assert.NotNil(t, result.Files) && assert.Len(t, result.Files, 1) {
		assert.Equal(t, "test.txt", result.Files[0].Filename, "Le fichier doit être correctement identifié")
	}
}
