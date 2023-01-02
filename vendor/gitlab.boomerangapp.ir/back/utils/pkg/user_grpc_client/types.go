package user_grpc_client

import "gitlab.boomerangapp.ir/back/utils/pkg/user_grpc_client/pb"

type UserProfile struct {
	Name       string
	Family     string
	Mobile     string
	NationalId string
	Email      string
	UserId     int64
}

func (up *UserProfile) resToUserProfile(res *pb.GetUserByIdResponse)  {
	up.Name = res.GetName()
	up.Family = res.GetFamily()
	up.Mobile = res.GetMobile()
	up.NationalId = res.GetNationalId()
	up.Email = res.GetEmail()
	up.UserId = res.GetUserId()
}