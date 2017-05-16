include "shared.thrift"

service Store extends shared.SharedService{
	void order(1:string product, 2:i32 amount),
	i32 getPrice(1:string product)
}
