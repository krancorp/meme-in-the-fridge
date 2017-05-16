include "shared.thrift"

service Store extends shared.SharedService{
	void order(1:string product, 2:i32 amount),
	double getPrice(1:string product)
}
