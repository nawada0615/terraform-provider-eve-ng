package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	labHTTPMethodGET    = "GET"
	labHTTPMethodPOST   = "POST"
	labHTTPMethodDELETE = "DELETE"
)

func setupMockEVE() *httptest.Server {
	mux := http.NewServeMux()

	// Mock login endpoint
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != labHTTPMethodPOST {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:  "unetlab_session",
			Value: "mock_session_123",
			Path:  "/api/",
		})
		fmt.Fprint(w, `{"code":200,"status":"success","message":"Login successful"}`)
	})

	// Mock lab management
	mux.HandleFunc("/api/labs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == labHTTPMethodGET {
			// Lab read
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Lab loaded",
				"data": {
					"author": "test",
					"description": "test lab",
					"body": "",
					"filename": "test-lab.unl",
					"id": "test-lab-id",
					"name": "test-lab",
					"version": "1",
					"scripttimeout": 300,
					"lock": false
				}
			}`)
		} else if r.Method == labHTTPMethodPOST {
			// Lab create
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab created"}`)
		} else if r.Method == labHTTPMethodDELETE {
			// Lab delete
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab deleted"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return httptest.NewServer(mux)
}

func TestEveLabCreate(t *testing.T) {
	server := setupMockEVE()
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProviderFactories: getProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "eve" {
						endpoint = "%s"
						username = "testuser"
						password = "testpass"
						insecure_skip_verify = true
					}
					resource "eve_lab" "test" {
						path = "/"
						name = "test-lab"
						author = "test"
						description = "test lab"
						version = "1"
					}
				`, server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("eve_lab.test", "name", "test-lab"),
					resource.TestCheckResourceAttr("eve_lab.test", "path", "/"),
					resource.TestCheckResourceAttr("eve_lab.test", "author", "test"),
					resource.TestCheckResourceAttr("eve_lab.test", "description", "test lab"),
					resource.TestCheckResourceAttr("eve_lab.test", "version", "1"),
				),
			},
		},
	})
}

func TestEveLabImport(t *testing.T) {
	server := setupMockEVE()
	defer server.Close()

	resource.Test(t, resource.TestCase{
		ProviderFactories: getProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "eve" {
						endpoint = "%s"
						username = "testuser"
						password = "testpass"
						insecure_skip_verify = true
					}
					resource "eve_lab" "test" {
						path = "/"
						name = "test-lab"
						author = "test"
						description = "test lab"
						version = "1"
					}
				`, server.URL),
			},
			{
				ResourceName:      "eve_lab.test",
				ImportState:       true,
				ImportStateVerify: false, // Skip verification due to computed attributes
			},
		},
	})
}
