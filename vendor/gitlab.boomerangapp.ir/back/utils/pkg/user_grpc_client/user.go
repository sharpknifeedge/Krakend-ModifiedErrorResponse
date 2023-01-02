package user_grpc_client

import (
	"context"
	utilsio "gitlab.boomerangapp.ir/back/utils/pkg/io"
	"gitlab.boomerangapp.ir/back/utils/pkg/user_grpc_client/pb"
	"google.golang.org/grpc"
	"time"
)

func GetUserById(address string, id uint) (*UserProfile, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer utilsio.SafeClose(conn)
	client := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	createFileRequest := pb.GetUserByIdRequest{
		UserId: int64(id),
	}
	res, err := client.GetUserById(ctx, &createFileRequest)
	if err != nil {
		return nil, err
	}
	up := new(UserProfile)
	up.resToUserProfile(res)
	return up, nil
}
