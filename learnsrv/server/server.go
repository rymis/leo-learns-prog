// Server part implementation for file storage
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/rymis/leo-learns-prog/learnsrv/rcs"
)

type Server struct {
	Root       string
	StaticRoot string

	Mux *http.ServeMux
}

type Context struct {
	User string
}

func New(root, static string) (*Server, error) {
	res := &Server{}
	res.Root = root
	res.StaticRoot = static
	res.Mux = http.NewServeMux()
	res.Mux.Handle("/", http.FileServer(http.Dir(static)))
	res.setupAPIHandlers()

	return res, nil
}

func (srv *Server) setupAPIHandlers() {
	srv.Mux.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		handleReq(w, r, srv.parseContext(r), srv.postEcho)
	})
	srv.Mux.HandleFunc("/api/commit", func(w http.ResponseWriter, r *http.Request) {
		handleReq(w, r, srv.parseContext(r), srv.postCommit)
	})
	srv.Mux.HandleFunc("/api/checkout", func(w http.ResponseWriter, r *http.Request) {
		handleReq(w, r, srv.parseContext(r), srv.postCheckout)
	})
	srv.Mux.HandleFunc("/api/versions", func(w http.ResponseWriter, r *http.Request) {
		handleReq(w, r, srv.parseContext(r), srv.postVersions)
	})
	srv.Mux.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		handleReq(w, r, srv.parseContext(r), srv.postList)
	})
}

func (srv *Server) parseContext(r *http.Request) *Context {
	return &Context{"leo"}
}

func (srv *Server) postEcho(ctx *Context, s *string) (string, error) {
	if s != nil {
		return *s, nil
	}
	return "", nil
}

type commitArgs struct {
	Content string `json:"content"`
	Comment string `json:"comment"`
	Name    string `json:"name"`
}

func (srv *Server) postCommit(ctx *Context, args *commitArgs) (string, error) {
	d, err := srv.getUserDir(ctx)
	if err != nil {
		return "", err
	}

	r, err := rcs.NewRCSFile(filepath.Join(d, args.Name))
	if err != nil {
		return "", err
	}

	v, err := r.Put(args.Content, args.Comment)

	return v, err
}

type checkoutArgs struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (srv *Server) postCheckout(ctx *Context, args *checkoutArgs) (string, error) {
	d, err := srv.getUserDir(ctx)
	if err != nil {
		return "", err
	}

	r, err := rcs.NewRCSFile(filepath.Join(d, args.Name))
	if err != nil {
		return "", err
	}

	var content string
	if args.Version == "" {
		content, err = r.Get()
	} else {
		content, err = r.GetVersion(args.Version)
	}

	return content, err
}

type versionsResult struct {
	Name     string            `json:"name"`
	Versions []rcs.VersionInfo `json:"versions"`
}

func (srv *Server) postVersions(ctx *Context, name *string) (versionsResult, error) {
	res := versionsResult{}

	if name == nil {
		return res, fmt.Errorf("Empty file name")
	}

	d, err := srv.getUserDir(ctx)
	if err != nil {
		return res, err
	}

	r, err := rcs.NewRCSFile(filepath.Join(d, *name))
	if err != nil {
		return res, err
	}

	res.Name = *name
	res.Versions, err = r.Versions()

	return res, err
}

func (srv *Server) postList(ctx *Context, mask *string) ([]string, error) {
	d, err := srv.getUserDir(ctx)
	if err != nil {
		return nil, err
	}

	files := rcs.ListFiles(d)
	if mask != nil {
		f := make([]string, 0, len(files))
		for _, fnm := range files {
			bnm := filepath.Base(fnm)
			if ok, err := path.Match(*mask, bnm); ok && err == nil {
				f = append(f, fnm)
			}
		}

		files = f
	}

	return files, nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.Mux.ServeHTTP(w, r)
}

func (srv *Server) getUserDir(ctx *Context) (string, error) {
	// TODO: check name correctness
	path := filepath.Join(srv.Root, "src", ctx.User)
	st, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return "", err
		}
		return path, nil
	}

	if !st.IsDir() {
		return "", fmt.Errorf("Not a directory: %v", err)
	}

	return path, nil
}

func errResponse(w http.ResponseWriter, err error) {
	type T struct {
		Error string `json:"error"`
	}

	res := &T{err.Error()}
	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}

func handleReq[T interface{}, R interface{}](w http.ResponseWriter, req *http.Request, ctx *Context, f func(ctx *Context, a *T) (R, error)) {
	w.Header().Add("content-type", "application/json")

	var args T
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&args)
	if err != nil {
		errResponse(w, err)
		return
	}

	type Res struct {
		Result R `json:"result"`
	}
	var res Res
	res.Result, err = f(ctx, &args)
	if err != nil {
		errResponse(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(&res)
}
