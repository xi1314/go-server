package admin_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAccession(t *testing.T) {
	r := admin.GetAccession()

	assert.Equal(t, accession.AdminList, r.Data)
}

func TestGetAccessionRouter(t *testing.T) {
	adminInfo, err := tester.LoginAdmin()

	assert.Nil(t, err)

	header := mocker.Header{
		"Authorization": token.JoinPrefixToken(adminInfo.Token),
	}

	r := tester.HttpAdmin.Get("/v1/admin/accession", nil, &header)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	var dataList []*accession.Accession

	assert.Nil(t, res.Decode(&dataList))

	assert.Equal(t, accession.AdminList, dataList)
}
