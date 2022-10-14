package main

import (
	"bufio"
	"context"
	"goods"
	"log"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	Adress = "localhost:9988"
)

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func main() {
	// 声明监听
	conn, err := grpc.Dial(Adress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	//客户端绑定连接
	client := goods.NewGoodsInfoClient(conn)
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name,"\n")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSuffix(desc,"\n")
	goodsNew := goods.Goods{
		Name: name,
		Desc: desc,
	}

	//添加一个商品
	res, errs := client.AddGoods(context.Background(), &goodsNew)
	if errs != nil {
		grpclog.Error(errs)
	}
	//返回结果
	log.Printf("%+v", res)

	//搜索一个商品
	goodsIdNew := goods.GoodsId{}
	goodsIdNew.Value = res.Value
	res2, errs2 := client.GetGoods(context.Background(), &goodsIdNew)

	if errs2 != nil {
		grpclog.Error(errs2)
	}
	//返回结果
	log.Printf("%v", res2)

}
