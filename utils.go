package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

func escape(s string) string {
	return s
}

func isErr(err error, code uint16) bool {
	return strings.HasPrefix(err.Error(), fmt.Sprintf("LDAP Result Code %d %q", code, ldap.LDAPResultCodeMap[code]))
}

func (s *LDAPManagerServer) findGroup(groupName string, attributes []string) (*ldap.SearchResult, error) {
	return s.ldap.Search(ldap.NewSearchRequest(
		s.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escape(groupName)),
		attributes,
		[]ldap.Control{},
	))
}

func (s *LDAPManagerServer) updateLastID(cn string, newID int) error {
	modifyRequest := ldap.NewModifyRequest(
		fmt.Sprintf("cn=%s,%s", cn, s.BaseDN),
		[]ldap.Control{},
	)
	modifyRequest.Replace("serialNumber", []string{strconv.Itoa(newID)})
	log.Debug(modifyRequest)
	if err := s.ldap.Modify(modifyRequest); err != nil {
		return fmt.Errorf("failed to update cn=%s: %v", cn, err)
	}
	log.Debugf("updated cn=%s with %d", cn, newID)
	return nil
}

/*
func serve() error {

	// The username and password we want to check
	// username := "someuser"
	password := "userpassword"

	bindusername := "uid=billy,ou=users,dc=example,dc=org"
	bindpassword := "hallo"

	bindusername = "cn=admin,dc=example,dc=org"
	bindpassword = "admin"

	// First bind with a read only user
	log.Println("Bind")
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username
	// SRCH base="ou=groups,dc=example,dc=org" scope=2 deref=0 filter="(cn=admins)"

	log.Println("Search")
	searchRequest := ldap.NewSearchRequest(
		"ou=groups,dc=example,dc=org",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(cn=admins)",
		[]string{"dn"},
		nil,
	)


	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=org",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn"},
		nil,
	)


	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(sr.Entries)
	if len(sr.Entries) != 1 {
		log.Fatal("User does not exist or too many entries returned")
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	log.Println("Bind as user")
	err = l.Bind(userdn, password)
	if err != nil {
		log.Fatal(err)
	}

	// Rebind as the read only user for any further queries
	log.Println("Rebind read only")
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
*/

/*
type query struct {
	DownloadToken     string `form:"download-token" json:"download-token"`
}

// DownloadHandler ...
func (s *LDAPManagerServer) DownloadHandler(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "DownloadHandler")

	// Get all download params
	licenseParam := c.Params.ByName("license")
	typeParam := c.Params.ByName("type")
	idParam := c.Params.ByName("id")

	var downloadToken string
	if c.GetHeader("download-token") != "" {
		downloadToken = c.GetHeader("download-token")
	} else {
		// URL query param as with the emails
		var q query
		if err := c.Bind(&q); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		downloadToken = q.DownloadToken
	}

	// Validate license and stock type
	lc, ok := licensepb.License_Type_value[strings.ToUpper(licenseParam)]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	typ, ok := stockpb.Stock_Type_value[strings.ToUpper(typeParam)]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	span.LogFields(
		slog.String("licenseParam", licenseParam),
		slog.String("typeParam", typeParam),
		slog.String("idParam", idParam),
		slog.String("downloadToken", downloadToken),
		slog.Int32("license", lc),
		slog.Int32("type", typ),
	)

	requestedID := &stockpb.Stock_ID{
		Type: stockpb.Stock_Type(typ),
		Id: int32(id),
	}

	// Validate the download here
	result, err := s.downloadService.ValidateDownload(ctx, &downloadmgmtpb.ValidateDownloadRequest{
		Token: downloadToken,
		Id: requestedID,
		LicenseType: licensepb.License_Type(lc),
	})
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	remaining := result.GetRemainingDownloads()
	if remaining > 0 {
		// go ahead with the download
		bucket := storage.BucketForLicense(licensepb.License_Type(lc))
		objName := storage.FormatStockID(requestedID, "")
		object, err := s.BaseService.OS.Client.GetObject(ctx, bucket, objName, minio.GetObjectOptions{})
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		info, err := object.Stat()
		if err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		// Decrement remaining downloads
		if _, err := s.downloadService.SetRemainingDownloads(ctx, &downloadmgmtpb.SetRemainingDownloadsRequest{
			RemainingDownloads: remaining - 1,
		}); err != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		contentLength := info.Size
		contentType := info.ContentType
		span.LogFields(
			slog.Int64("contentLength", contentLength),
			slog.String("contentType", contentType),
		)

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
			"Content-Type": contentType,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, object, extraHeaders)
	} else {
		c.Status(http.StatusGone)
	}
}
*/
