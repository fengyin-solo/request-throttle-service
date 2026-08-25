package store
import("testing";"ratelimiter/internal/model";"ratelimiter/pkg/idgen")
func TestClientLookupDoesNotExposeMutableStoreState(t *testing.T){s:=NewMemoryStore();c:=&model.Client{ID:idgen.Hex(),AppID:"app",Name:"before"};_ = s.CreateClient(c);got,_:=s.GetClient(c.ID);got.Name="changed";again,_:=s.GetClient(c.ID);if again.Name!="before"{t.Fatalf("lookup exposed mutable client: %q",again.Name)}}
