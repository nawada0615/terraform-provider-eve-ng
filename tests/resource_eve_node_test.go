package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	nodeHTTPMethodGET    = "GET"
	nodeHTTPMethodPOST   = "POST"
	nodeHTTPMethodDELETE = "DELETE"
)

func setupMockEVEForNode() *httptest.Server {
	mux := http.NewServeMux()

	// Mock login endpoint
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != nodeHTTPMethodPOST {
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

		if r.Method == nodeHTTPMethodPOST {
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

		if r.Method == nodeHTTPMethodGET {
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
		} else if r.Method == nodeHTTPMethodPOST {
			// Lab create
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab created"}`)
		} else if r.Method == nodeHTTPMethodDELETE {
			// Lab delete
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Lab deleted"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Mock node management
	mux.HandleFunc("/api/labs/test-lab.unl/nodes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == nodeHTTPMethodGET {
			// Node list
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Nodes listed",
				"data": {
					"1": {
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
						"cpulimit": false,
						"uuid": "test-uuid-123",
						"firstmac": "00:50:00:00:01:00",
						"qemu_version": "2.4.0",
						"qemu_arch": "x86_64",
						"qemu_nic": "virtio-net-pci",
						"qemu_options": "-enable-kvm",
						"status": 0
					}
				}
			}`)
		} else if r.Method == nodeHTTPMethodPOST {
			// Node create
			fmt.Fprint(w, `{
				"code": 201,
				"status": "success",
				"message": "Node created",
				"data": {
					"id": 1
				}
			}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Mock individual node management
	mux.HandleFunc("/api/labs/test-lab.unl/nodes/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == nodeHTTPMethodGET {
			// Node read
			fmt.Fprint(w, `{
				"code": 200,
				"status": "success",
				"message": "Node retrieved",
				"data": {
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
				}
			}`)
		} else if r.Method == "PUT" {
			// Node update
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Node updated"}`)
		} else if r.Method == nodeHTTPMethodDELETE {
			// Node delete
			fmt.Fprint(w, `{"code":200,"status":"success","message":"Node deleted"}`)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return httptest.NewServer(mux)
}

func TestEveNodeCreate(t *testing.T) {
	server := setupMockEVEForNode()

	nodeConfig := `resource "eve_node" "test" {
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
	}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("eve_node.test", "name", "test-node"),
		resource.TestCheckResourceAttr("eve_node.test", "type", "qemu"),
		resource.TestCheckResourceAttr("eve_node.test", "template", "linux"),
		resource.TestCheckResourceAttr("eve_node.test", "image", "linux-ubuntu-22.04"),
		resource.TestCheckResourceAttr("eve_node.test", "cpu", "1"),
		resource.TestCheckResourceAttr("eve_node.test", "ram", "1024"),
		resource.TestCheckResourceAttr("eve_node.test", "ethernet", "4"),
	}

	runResourceTest(t, server, nodeConfig, checks)
}
