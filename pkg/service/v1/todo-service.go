package v1

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.smartmachine.io/go-grpc-api/pkg/api/v1"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	db *gorm.DB
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(db *gorm.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db: db}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Create new todo task
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	_, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	orm, err := req.ToDo.ToORM(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert to orm representation: " + err.Error())
	}

	err = s.db.Create(&orm).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to insert, internal error: %v", err)
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  orm.Id,
	}, nil
}

// Read todo task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	var orm v1.ToDoORM
	err := s.db.First(&orm, req.Id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "record not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "error reading record: %v", err)
	}

	td, err := orm.ToPB(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to convert to orm representation: %v", err)
	}

	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &td,
	}, nil

}

// Update todo task
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	orm, err := req.ToDo.ToORM(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to convert to orm representation: " + err.Error())
	}

	err = s.db.Save(orm).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Errorf(codes.NotFound, "record not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "error updating record: %v", err)
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: 1,
	}, nil
}

// Delete todo task
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	db := s.db.Delete(&v1.ToDoORM{Id: req.Id})
	err := db.Error
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, status.Errorf(codes.NotFound, "unable to delete, record not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "unable to delete, internal error: %v", err)
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: db.RowsAffected,
	}, nil
}

// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	var users []*v1.ToDoORM
	var count int

	err := s.db.Find(&users).Count(&count).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, status.Errorf(codes.NotFound, "unable to read records, not found: %v", err)

	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to read records, internal error: %v", err)
	}

	list := []*v1.ToDo{}
	for _, todo := range users {

		orm, err := todo.ToPB(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "unable to convert to orm representation: " + err.Error())
		}

		list = append(list, &orm)
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: list,
	}, nil
}