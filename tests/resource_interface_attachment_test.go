package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	interfaceHTTPMethodGET    = "GET"
	interfaceHTTPMethodPOST   = "POST"
	interfaceHTTPMethodPUT    = "PUT"
	interfaceHTTPMethodDELETE = "DELETE"
)

func setupMockEVEForInterfaceAttachment() *httptest.Server {
	mux := http.NewServeMux()
	setupLoginEndpoint(mux)
	setupLabEndpoints(mux)
	setupNetworkEndpoints(mux)
	setupNodeEndpoints(mux)
	setupInterfaceEndpoints(mux)
	return httptest.NewServer(mux)
}

func setupLoginEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != interfaceHTTPMethodPOST {
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
}

func setupLabEndpoints(mux *http.ServeMux) {
	setupLabCreationEndpoint(mux)
	setupLabManagementEndpoint(mux)
}

func setupLabCreationEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("/api/labs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == interfaceHTTPMethodPOST {
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
}

func setupLabManagementEndpoint(mux *http.ServeMux) {
	mux.HandleFunc("/api/labs/test-lab.unl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == interfaceHTTPMethodGET {
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
		} else if r.Method == interfaceHTTPMethodPOST {
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab created"}`)
		} else if r.Method == interfaceHTTPMethodDELETE {
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab deleted"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// setupResourceEndpoints creates common CRUD endpoints for a resource type
func setupResourceEndpoints(mux *http.ServeMux, resourceType, resourceName, dataTemplate string) {
	// Resource list and create
	mux.HandleFunc("/api/labs/test-lab.unl/"+resourceType, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == interfaceHTTPMethodGET {
			fmt.Fprintf(w, `{
				"code": 200,
				"status": "success",
				"message": "%ss listed",
				"data": {
					"1": %s
				}
			}`, resourceName, dataTemplate)
		} else if r.Method == interfaceHTTPMethodPOST {
			fmt.Fprintf(w, `{
				"code": 201,
				"status": "success",
				"message": "%s created",
				"data": {
					"id": 1
				}
			}`, resourceName)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Individual resource management
	mux.HandleFunc("/api/labs/test-lab.unl/"+resourceType+"/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == interfaceHTTPMethodGET {
			fmt.Fprintf(w, `{
				"code": 200,
				"status": "success",
				"message": "%s retrieved",
				"data": %s
			}`, resourceName, dataTemplate)
		} else if r.Method == interfaceHTTPMethodPUT {
			fmt.Fprintf(w, `{"code":200,"status":"success","message":"%s updated"}`, resourceName)
		} else if r.Method == interfaceHTTPMethodDELETE {
			fmt.Fprintf(w, `{"code":200,"status":"success","message":"%s deleted"}`, resourceName)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func setupNetworkEndpoints(mux *http.ServeMux) {
	setupResourceEndpoints(mux, "networks", "Network", `{
		"name": "test-net",
		"type": "bridge",
		"icon": "cloud.png",
		"top": 200,
		"left": 200,
		"visibility": "1",
		"node_count": 0
	}`)
}

func setupNodeEndpoints(mux *http.ServeMux) {
	setupResourceEndpoints(mux, "nodes", "Node", `{
		"name": "test-node",
		"type": "qemu",
		"template": "linux",
		"image": "linux-ubuntu-22.04",
		"icon": "Router-2D-Gen-White-S.svg",
		"top": 100,
		"left": 100,
		"delay": 0,
		"cpu": 1,
		"ram": 1024,
		"ethernet": 4,
		"serial": 0,
		"status": 0
	}`)
}

func setupInterfaceEndpoints(mux *http.ServeMux) {
	// Interface attachment management
	mux.HandleFunc("/api/labs/test-lab.unl/nodes/1/interfaces", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == interfaceHTTPMethodGET {
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Interfaces listed",
				"data": {
					"0": {
						"name": "eth0",
						"type": "ethernet",
						"target": "network:1"
					}
				}
			}`)
		} else if r.Method == interfaceHTTPMethodPOST {
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Interface attached"}`)
		} else if r.Method == interfaceHTTPMethodPUT {
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Interface updated"}`)
		} else if r.Method == interfaceHTTPMethodDELETE {
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Interface detached"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func TestEveInterfaceAttachmentCreate(t *testing.T) {
	server := setupMockEVEForInterfaceAttachment()
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
					resource "eve_network" "test" {
						lab_file = eve_lab.test.file
						name = "test-net"
						type = "bridge"
						icon = "cloud.png"
						top = 200
						left = 200
						visibility = "1"
					}
					resource "eve_node" "test" {
						lab_file = eve_lab.test.file
						name = "test-node"
						type = "qemu"
						template = "linux"
						image = "linux-ubuntu-22.04"
						icon = "Router-2D-Gen-White-S.svg"
						top = 100
						left = 100
						cpu = 1
						ram = 1024
						ethernet = 4
						desired_state = "stopped"
					}
					resource "eve_interface_attachment" "test" {
						lab_file        = eve_lab.test.file
						node_id         = tonumber(split(":node:", eve_node.test.id)[1])
						interface_index = 0
						target          = "network:${tonumber(split(":network:", eve_network.test.id)[1])}"
					}
				`, server.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("eve_interface_attachment.test", "lab_file", "/test-lab.unl"),
					resource.TestCheckResourceAttr("eve_interface_attachment.test", "node_id", "1"),
					resource.TestCheckResourceAttr("eve_interface_attachment.test", "interface_index", "0"),
					resource.TestCheckResourceAttr("eve_interface_attachment.test", "target", "network:1"),
				),
			},
		},
	})
}
