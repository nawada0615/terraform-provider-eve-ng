package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	networkHTTPMethodGET    = "GET"
	networkHTTPMethodPOST   = "POST"
	networkHTTPMethodDELETE = "DELETE"
)

func setupMockEVEForNetwork() *httptest.Server {
	mux := http.NewServeMux()

	// Mock login endpoint
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != networkHTTPMethodPOST {
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

	// Mock lab creation endpoint
	mux.HandleFunc("/api/labs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == networkHTTPMethodPOST {
			// Lab create
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Lab created",
				"data": {
					"filename": "test-lab.unl"
				}
			}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Mock lab management
	mux.HandleFunc("/api/labs/test-lab.unl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == networkHTTPMethodGET {
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
		} else if r.Method == networkHTTPMethodPOST {
			// Lab create
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab created"}`)
		} else if r.Method == networkHTTPMethodDELETE {
			// Lab delete
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab deleted"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Mock network management
	mux.HandleFunc("/api/labs/test-lab.unl/networks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == networkHTTPMethodGET {
			// Network list
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Networks listed",
				"data": {
					"1": {
						"name": "test-net",
						"type": "bridge",
						"icon": "cloud.png",
						"top": 200,
						"left": 200,
						"visibility": "1",
						"node_count": 0
					}
				}
			}`)
		} else if r.Method == networkHTTPMethodPOST {
			// Network create
			fmt.Fprint(w, `{
				"code": 201,
				"status": "success",
				"message": "Network created",
				"data": {
					"id": 1
				}
			}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Mock individual network management
	mux.HandleFunc("/api/labs/test-lab.unl/networks/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == networkHTTPMethodGET {
			// Network read
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Network retrieved",
				"data": {
					"name": "test-net",
					"type": "bridge",
					"icon": "cloud.png",
					"top": 200,
					"left": 200,
					"visibility": "1",
					"node_count": 0
				}
			}`)
		} else if r.Method == interfaceHTTPMethodPUT {
			// Network update
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Network updated"}`)
		} else if r.Method == networkHTTPMethodDELETE {
			// Network delete
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Network deleted"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return httptest.NewServer(mux)
}

func TestEveNetworkCreate(t *testing.T) {
	server := setupMockEVEForNetwork()

	networkConfig := `resource "eve_network" "test" {
		lab_file = eve_lab.test.file
		name = "test-net"
		type = "bridge"
		icon = "cloud.png"
		top = 200
		left = 200
		visibility = "1"
	}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("eve_network.test", "name", "test-net"),
		resource.TestCheckResourceAttr("eve_network.test", "type", "bridge"),
		resource.TestCheckResourceAttr("eve_network.test", "icon", "cloud.png"),
		resource.TestCheckResourceAttr("eve_network.test", "top", "200"),
		resource.TestCheckResourceAttr("eve_network.test", "left", "200"),
		resource.TestCheckResourceAttr("eve_network.test", "visibility", "1"),
	}

	runResourceTest(t, server, networkConfig, checks)
}
