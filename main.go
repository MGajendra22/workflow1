package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"log"
	"encoding/json"

)

type task struct {
	Id     int     `json:"id"`
	Desc   string  `json:"desc"`
	Status bool    `json:"status"`
}

type tasks struct {
	cont []task
}


type handle struct {
	ptr *tasks
}

// var t1 tasks

var temp_id = idgen()

func idgen() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

func GetTaskById(id int, h *handle) *task {

	for _, val := range h.ptr.cont {

		if val.Id == id {
			return &val
		}
	}
	return nil
}

func (h* handle) handleAll(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case "GET" : {
		_, err := w.Write([]byte("List of all tasks :- \n"))
		if err != nil {
			fmt.Printf("Error")
		}
		for _, val := range h.ptr.cont {
			body,err:=json.Marshal(val)
			if err != nil {
            panic(err)
            }
			var t1 task
			err=json.Unmarshal(body,&t1)
			if err != nil {
            panic(err)
            }
			 _,err=w.Write([]byte(body))
			 if err != nil {
            panic(err)
            }

		}
	}

	case "POST" :{
		bodyBy, err := io.ReadAll(r.Body)
		fmt.Println(string(bodyBy))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("invalid body"))
			if err != nil {
				fmt.Printf("Error")
			}
			return
		}

		desc := string(bodyBy)
		if desc == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			t := task{temp_id(), desc, false}
			body,err:=json.Marshal(t)
			if err != nil {
            panic(err)
            }
			var t1 task
			err=json.Unmarshal(body,&t1)
			if err != nil {
            panic(err)
            }
			h.ptr.cont = append(h.ptr.cont, t1)
			w.WriteHeader(http.StatusCreated)
			_,err=w.Write([]byte(body))
			if err != nil {
            panic(err)
            }
			
		}
	}

	case "PUT" :{

		idstr := r.PathValue("id")
		id, err := strconv.Atoi(idstr)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if h.ptr.cont[id].Status {
			_, err := w.Write([]byte("Task is already Completed"))
			if err != nil {
				fmt.Printf("Error")
			}
		} else {

			h.ptr.cont[id].Status = true
			str := fmt.Sprintf("Task %d is completed ", id)
			_, err := w.Write([]byte(str))
			if err != nil {
				fmt.Printf("Error")
			}

		}

	}

	}
}




func (h *handle) handlePendingTasks(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("List of Pending tasks :- \n"))
		if err != nil {
			fmt.Printf("Error")
		}
		for _, val := range h.ptr.cont {
			if !val.Status {
				body,err:=json.Marshal(val)
			if err != nil {
            panic(err)
            }
			var t1 task
			err=json.Unmarshal(body,&t1)
			if err != nil {
            panic(err)
            }
		    _,err=w.Write([]byte(body))
				if err != nil {
					fmt.Printf("Error")
				}
			}
		}

	}
}

func main() {

	a1 := task{temp_id(), " Grocery ", false}
	a2 := task{temp_id(), " Shopping ", false}

	var h handle

	h.ptr = &tasks{}

	h.ptr.cont = append(h.ptr.cont, a1, a2)

	http.HandleFunc("/task", h.handleAll)
	http.HandleFunc("/task/{id}", h.handleAll)
	http.HandleFunc("/pending", h.handlePendingTasks)


	fmt.Println("Server Started running on http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
