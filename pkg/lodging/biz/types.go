package lodging_biz

type Location struct {
	Type        string    `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Coordinates []float64 `protobuf:"fixed64,2,rep,packed,name=coordinates,proto3" json:"coordinates,omitempty"`
}

type VendorType struct {
	Name          string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Img           string   `protobuf:"bytes,2,opt,name=img,proto3" json:"img,omitempty"`
	MemberSince   string   `protobuf:"bytes,3,opt,name=member_since,json=memberSince,proto3" json:"member_since,omitempty"`
	Languages     []string `protobuf:"bytes,4,rep,name=languages,proto3" json:"languages,omitempty"`
	ResponseRate  int32    `protobuf:"varint,5,opt,name=response_rate,json=responseRate,proto3" json:"response_rate,omitempty"`
	ResponseTime  string   `protobuf:"bytes,6,opt,name=response_time,json=responseTime,proto3" json:"response_time,omitempty"`
	Location      string   `protobuf:"bytes,7,opt,name=location,proto3" json:"location,omitempty"`
	BoatName      string   `protobuf:"bytes,8,opt,name=boat_name,json=boatName,proto3" json:"boat_name,omitempty"`
	BoatGuests    int32    `protobuf:"varint,9,opt,name=boat_guests,json=boatGuests,proto3" json:"boat_guests,omitempty"`
	BoatCabins    int32    `protobuf:"varint,10,opt,name=boat_cabins,json=boatCabins,proto3" json:"boat_cabins,omitempty"`
	BoatBathrooms int32    `protobuf:"varint,11,opt,name=boat_bathrooms,json=boatBathrooms,proto3" json:"boat_bathrooms,omitempty"`
	TotalReview   int32    `protobuf:"varint,12,opt,name=total_review,json=totalReview,proto3" json:"total_review,omitempty"`
}

type EquipmentType struct {
	Img  string `protobuf:"bytes,1,opt,name=img,proto3" json:"img,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

type SpecificationType struct {
	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Details string `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
}

type ReviewType struct {
	Avatar   string `protobuf:"bytes,1,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Date     string `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	Rating   int32  `protobuf:"varint,4,opt,name=rating,proto3" json:"rating,omitempty"`
	Location string `protobuf:"bytes,5,opt,name=location,proto3" json:"location,omitempty"`
	Review   string `protobuf:"bytes,6,opt,name=review,proto3" json:"review,omitempty"`
}

type ReviewBarType struct {
	Count   int32   `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	Percent float32 `protobuf:"fixed32,2,opt,name=percent,proto3" json:"percent,omitempty"`
}

type ReviewStatsType struct {
	TotalReviews  int32           `protobuf:"varint,1,opt,name=total_reviews,json=totalReviews,proto3" json:"total_reviews,omitempty"`
	AverageRating float32         `protobuf:"fixed32,2,opt,name=average_rating,json=averageRating,proto3" json:"average_rating,omitempty"`
	Stars         []*ReviewBarType `protobuf:"bytes,3,rep,name=stars,proto3" json:"stars,omitempty"`
}