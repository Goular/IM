#如何安全接入这个聊天主页面
通过登录api `user/login` 获得id和token
返回json到前端,
前端拼接url
/chat/index.shtml?id=1&token=123
通过location.href 跳转

#如何添加/显示好友 添加/显示群
```cgo
/contact/addfriend 自动添加好友,参数userid,dstid
//
用户10000添加好友10086,往contact表中添加俩条记录

//一条ownerid =10000,dstobj=10086 
//一条ownerid =10086,dstobj=10000

/contact/loadfriend 显示全部好友,参数userid

/contact/createcommunity 建群,头像pic,名称name,备注memo,创建者userid
/contact/loadcommunity 显示全部群 参数userid
//加群逻辑特殊一点
/contact/joincommunity 加群,参数userid,dstid

```

##创建模型(实体)
```cgo
const (
		CONCAT_CATE_USER = 0x01  //用户
	    CONCAT_CATE_COMUNITY = 0x02 //群组
	)
type Contact struct {
	Id         int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Ownerid       int64	`xorm:"bigint(20)" form:"ownerid" json:"ownerid"`   // 什么角色
	Dstobj       int64	`xorm:"bigint(20)" form:"dstobj" json:"dstobj"`   // 什么角色
	Cate      int	`xorm:"int(11)" form:"cate" json:"cate"`   // 什么角色
	Memo    string	`xorm:"varchar(120)" form:"memo" json:"memo"`   // 什么角色
	Createat   time.Time	`xorm:"datetime" form:"createat" json:"createat"`   // 什么角色
}

//同步表结构
DbEngin.Sync2(new(model.Contact))
```
##创建控制器ctrl
```cgo
func Addfriend(w http.ResponseWriter, req *http.Request) {
	//request.ParseForm()
    //mobile := request.PostForm.Get("mobile")
    //passwd := request.PostForm.Get("passwd")
	//str->int
	//
	var arg args.ContactArg
	//对象绑定
	util.Bind(req,&arg)
	//
	err := contactService.AddFriend(arg.Userid,arg.Dstid)
	
	if err!=nil{
		util.RespFail(w,err.Error())
	}else{
		util.RespOk(w,msgs)
	}
}
```
##配置路由
http.HandleFunc("/contact/addfriend", ctrl.Addfriend)
##配置service
```cgo
//自动添加好友
func (service *ContactService) AddFriend(
	userid,//用户id 10086,
	dstid int64 ) error{
	
	//判断是否已经存在
	//如果存在记录说明已经是好友了不加
	if tmp.Id>0{
		return errors.New("该用户已经被添加过啦")
	}
	//启动事务
	session := DbEngin.NewSession();
	session.Begin()
	//
	_,e2 := session.InsertOne()
	//
	_,e3 := session.InsertOne()
	//
	if  e2==nil && e3==nil{
		session.Commit()
		//
	}else{
		session.Rollback()
		//
	}
}
```

##前端js
```javascript
addfriend:function(){
    
    //弹窗提示用户输入
    // mui.prompt()
}
_addfriend:function(){
   //网络请求
}
```



















```cgo

var contactService service.ContactService
func LoadFriend(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	users := contactService.SearchFriend(arg.Userid)
	util.RespOkList(w,users,len(users))
}

func Loadcommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg
	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	comunitys := contactService.SearchComunity(arg.Userid)
	util.RespOkList(w,comunitys,len(comunitys))
}
func JoinCommunity(w http.ResponseWriter, req *http.Request){
	var arg args.ContactArg

	//如果这个用的上,那么可以直接
	util.Bind(req,&arg)
	err := contactService.JoinCommunity(arg.Userid,arg.Dstid);
	if err!=nil{
		util.RespFail(w,err.Error())
	}else {
		util.RespOk(w,"")
	}
}


http.HandleFunc("/contact/addfriend", ctrl.Loadcommunity)
http.HandleFunc("/contact/community", ctrl.Loadcommunity)
http.HandleFunc("/contact/friend", ctrl.LoadFriend)

```






















##创建服务层service
```cgo
//搜索群组
func (service *ContactService) SearchComunity(userId int64) ([]model.Community){
	 conconts := make([]model.Contact,0)
	 comIds :=make([]int64,0)

	 DBengin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUBITY).Find(&conconts)
     for _,v := range conconts{
		 comIds = append(comIds,v.Dstobj);
	 }
     coms := make([]model.Community,0)
     if len(comIds)== 0{
     	return coms
	 }
	DBengin.In("id",comIds).Find(&coms)
	return coms
}

//加群
func (service *ContactService) JoinCommunity(userId,comId int64) error{
	cot := model.Contact{
		Ownerid:userId,
		Dstobj:comId,
		Cate:model.CONCAT_CATE_COMUBITY,
	}
	DBengin.Get(&cot)
	if(cot.Id==0){
		cot.Createat = time.Now()
		_,err := DBengin.InsertOne(cot)
		return err
	}else{
		return nil
	}


}
//建群
func (service *ContactService) CreateCommunity(comm model.Community) (ret model.Community,err error){
	if len(comm.Name)==0{
		err = errors.New("缺少群名称")
		return ret,err
	}
	if comm.Ownerid==0{
		err = errors.New("请先登录")
		return ret,err
	}
	com := model.Community{
		Ownerid:comm.Ownerid,
	}
	num,err := DBengin.Count(&com)

	if(num>5){
		err = errors.New("一个用户最多只能创见5个群")
		return com,err
	}else{
		comm.Createat=time.Now()
		session := DBengin.NewSession()
		session.Begin()
		_,err = session.InsertOne(&comm)
		if err!=nil{
			session.Rollback();
			return com,err
		}
		_,err =session.InsertOne(
			model.Contact{
				Ownerid:comm.Ownerid,
				Dstobj:comm.Id,
				Cate:model.CONCAT_CATE_COMUBITY,
				Createat:time.Now(),
			})
		if err!=nil{
			session.Rollback();
		}else{
			session.Commit()
		}
		return com,err
	}
}
//加好友
func (service *ContactService) AddFriend(userId,dstId int64) error{
	cot := model.Contact{
		Ownerid:userId,
		Dstobj:dstId,
		Cate:model.CONCAT_CATE_USER,
	}
	DBengin.Get(&cot)
	if(cot.Id==0){
		cot.Createat = time.Now()
		_,err := DBengin.InsertOne(cot)
		return err
	}else{
		return nil
	}

}
//查找好友
func (service *ContactService) SearchFriend(userId int64) ([]model.User){
	conconts := make([]model.Contact,0)
	objIds :=make([]int64,0)
	DBengin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_USER).Find(&conconts)
	for _,v := range conconts{
		objIds = append(objIds,v.Dstobj);
	}
	coms := make([]model.User,0)
	if len(objIds)== 0{
		return coms
	}
	DBengin.In("id",objIds).Find(&coms)
	return coms
}
//自动添加好友
func (service *ContactService) AddFriendAuto(userid,dstid int64){
	//事务的使用
	session := DBengin.NewSession();
	session.Begin()
	
	_,e2 := session.Insert(model.Contact{
		Ownerid:userid,
		Dstobj:dstid,
		Cate:model.CONCAT_CATE_USER,
		Createat:time.Now(),
	})
	_,e3 := session.Insert(model.Contact{
		Ownerid:dstid,
		Dstobj:userid,
		Cate:model.CONCAT_CATE_USER,
		Createat:time.Now(),
	})
	if  e2==nil && e3==nil{
		session.Commit()
	}else{
		session.Rollback()
	}
}

```
##前端调用实现
```javascript
//核心工具包解析


```

```javascript
//解析query,/chat/index.shtml?id=1&token=x
//token = parseQuery("token")
function parseQuery (name){	
  var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
  var r = window.location.search.substr(1).match(reg);  //匹配目标参数
  if (r != null) return decodeURI(unescape(r[2])); 
  return null; //返回参数值
 }
 //获取用户ID用户ID
 function userId(){
    var id = parseQuery("id")
    if (id==null){
        return 0
    }else{
        return parseInt(id)
    }
 }
 //核心网络请求函数
 function post(uri,data,fn){
                
            var xhr = new XMLHttpRequest();
            xhr.open("POST",url, true);
            // 添加http头，发送信息至服务器时内容编码类型
            xhr.setRequestHeader(
            	"Content-Type",
    			"application/x-www-form-urlencoded"
    		);
            xhr.onreadystatechange = function() {
                if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
                    //resolve(JSON.parse(xhr.responseText));
                    if (typeof fn=="function"){
                        fn(JSON.parse(xhr.responseText))    
                    }
                }
            };
            xhr.onerror = function(){
            	//reject({"code":-1,"msg":"服务器繁忙"})
            	if (typeof fn=="function"){
            	    fn({"code":-1,"msg":"服务器繁忙"})
            	}
    		}
            var _data=[];
            for(var i in data){
                _data.push( i +"=" + encodeURI(data[i]));
            }
        	xhr.send(_data.join("&"));
        
 }
 //网络请求函数promis版本
 function post(uri,data,fn){
     return new Promise(function (resolve, reject) {
             var xhr = new XMLHttpRequest();
             xhr.open("POST",url, true);
             // 添加http头，发送信息至服务器时内容编码类型
             xhr.setRequestHeader(
             	"Content-Type",
     			"application/x-www-form-urlencoded"
     		);
             xhr.onreadystatechange = function() {
                 if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
                     resolve(JSON.parse(xhr.responseText));
                     if (typeof fn=="function"){
                         fn(JSON.parse(xhr.responseText))    
                     }
                 }
             };
             xhr.onerror = function(){
             	reject({"code":-1,"msg":"服务器繁忙"})
             	if (typeof fn=="function"){
                  fn(JSON.parse(xhr.responseText))    
                }
     		}
             var _data=[];
             for(var i in data){
                 _data.push( i +"=" + encodeURI(data[i]));
             }
         	xhr.send(_data.join("&"));
         })
  }
  
 loadfriends:function(){
           var that = this;
           post("contact/friend",{userid:userId()},
           function(res){
               that.friends = res.rows ||[];
               var usermap = usermap;
               for(var i in res.rows){
                   usermap[res.rows[i].id]=res.rows[i];
               }
               this.usermap = usermap;
           }.bind(this))
       },
       loadcommunitys:function(){
           var that = this;
           post("contact/community",{userid:userId()},function(res){
               that.communitys = res.rows ||[];
           })
       },
       addfriend:function(){
           var that = this;
           mui.prompt('','请输入好友ID','加好友',['取消','确认'],function (e) {
               if (e.index == 1) {
                   if (isNaN(e.value) || e.value <= 0) {
                       mui.toast('格式错误');
                   }else{
                       //mui.toast(e.value);
                       that._addfriend(e.value)
                   }
               }else{
                   //mui.toast('您取消了入库');
               }
           },'div');
           document.querySelector('.mui-popup-input input').type = 'number';
       },
       _addfriend:function(dstobj){
           var user = userInfo();
 
           post("user/addfriend",{dstid:dstobj,userid: user.id,pic:user.avatar,content:user.nickname,memo: "请求加你为好友"},function(res){
               if(res.code==200){
                   mui.toast("添加成功");
                   that.loadfriends();
               }else{
                   mui.toast(res.msg);
               }
           })
       },
       _joincomunity:function(dstobj){
           var that = this;
           post("user/joincommunity",{dstid:dstobj},function(res){
               if(res.code==200){
                   mui.toast("添加成功");
 
                   that.loadcommunitys();
               }else{
                   mui.toast(res.msg);
               }
           })
       },
       joincomunity:function(){
           var that = this;
           mui.prompt('','请输入群号','加群',['取消','确认'],function (e) {
               if (e.index == 1) {
                   if (isNaN(e.value) || e.value <= 0) {
                       mui.toast('格式错误');
                   }else{
                       //mui.toast(e.value);
                       that._joincomunity(e.value)
                   }
               }else{
                   //mui.toast('您取消了入库');
               }
           },'div');
           document.querySelector('.mui-popup-input input').type = 'number';
       },
```


#设计可以无限扩张业务场景的消息通讯结构
```cgo
func recvproc(node *Node) {
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		//dispatch(data)
		fmt.Printf("recv<=%s",data)
	}
}
```
##原理
前端通过websocket发送`json格式的字符串`
用户2向用户3发送文字消息hello
```json5
{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}
```
里面携带
谁发的-userid
要发给谁-dstid
这个消息有什么用-cmd
消息怎么展示-media
消息内容是什么-(url,amout,pic,content等)
##核心数据结构
```cgo
type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"` //消息ID
	//谁发的
	Userid  int64  `json:"userid,omitempty" form:"userid"` //谁发的
	//什么业务
	Cmd     int    `json:"cmd,omitempty" form:"cmd"` //群聊还是私聊
	//发给谁
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`//对端用户ID/群ID
	//怎么展示
	Media   int    `json:"media,omitempty" form:"media"` //消息按照什么样式展示
	//内容是什么
	Content string `json:"content,omitempty" form:"content"` //消息的内容
	//图片是什么
	Pic     string `json:"pic,omitempty" form:"pic"` //预览图片
	//连接是什么
	Url     string `json:"url,omitempty" form:"url"` //服务的URL
	//简单描述
	Memo    string `json:"memo,omitempty" form:"memo"` //简单描述
	//其他的附加数据，语音长度/红包金额
	Amount  int    `json:"amount,omitempty" form:"amount"` //其他和数字相关的
}
const (
    //点对点单聊,dstid是用户ID
	CMD_SINGLE_MSG = 10
	//群聊消息,dstid是群id
	CMD_ROOM_MSG   = 11
	//心跳消息,不处理
	CMD_HEART      = 0
	
)
const (
    //文本样式
	MEDIA_TYPE_TEXT=1
	//新闻样式,类比图文消息
	MEDIA_TYPE_News=2
	//语音样式
	MEDIA_TYPE_VOICE=3
	//图片样式
	MEDIA_TYPE_IMG=4
	
	//红包样式
	MEDIA_TYPE_REDPACKAGR=5
	//emoj表情样式
	MEDIA_TYPE_EMOJ=6
	//超链接样式
	MEDIA_TYPE_LINK=7
	//视频样式
	MEDIA_TYPE_VIDEO=8
	//名片样式
	MEDIA_TYPE_CONCAT=9
	//其他自己定义,前端做相应解析即可
	MEDIA_TYPE_UDEF=100
)
/**
消息发送结构体,点对点单聊为例
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,
content:"hello"}

3、MEDIA_TYPE_VOICE,amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,
url:"http://www.a,com/dsturl.mp3",
amount:40}

4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,
url:"http://www.baidu.com/a/log.jpg"}


2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,
content:"标题",
pic:"http://www.baidu.com/a/log,jpg",
url:"http://www.a,com/dsturl",
"memo":"这是描述"}


5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}

7、MEDIA_TYPE_Link 7
{id:1,userid:2,dstid:3,cmd:10,media:7,
"url":"http://www.a.com/dsturl.html"
}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,
pic:"http://www.baidu.com/a/log,jpg",
url:"http://www.a,com/a.mp4"
}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,
"content":"10086",
"pic":"http://www.baidu.com/a/avatar,jpg",
"memo":"胡大力"}

*/
```
##从哪里接收数据?怎么处理这些数据呢?
```cgo
func recvproc(node *Node) {
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		fmt.Printf("recv<=%s",data)
		dispatch(data)
	}
}
func dispatch(data []byte){
    //todo 转成message对象
    
    //todo 根据cmd参数处理逻辑
    
    
    
    
    
    
    msg :=Message{}
    err := json.UnMarshal(data,&msg)
    if err!=nil{
        log.Printf(err.Error())
        return ;
    }
    switch msg.Cmd {
    	case CMD_SINGLE_MSG: //如果是单对单消息,直接将消息转发出去
    		//向某个用户发回去
    		fmt.Printf("c2cmsg %d=>%d\n%s\n",msg.Userid,msg.Dstid,string(tmp))
    		SendMsgToUser(msg.Userid, msg.Dstid, tmp)
    		//fmt.Println(msg)
    	case CMD_ROOM_MSG: //群聊消息,需要知道
    		fmt.Printf("c2gmsg %d=>%d\n%s\n",msg.Userid,msg.Dstid,string(tmp))
    		SendMsgToRoom(msg.Userid, msg.Dstid, tmp)
    	case CMD_HEART:
    	default:
    	    //啥也别做
    	    
    	}
    		
}
```



#5.4 实现发送文字、表情包等

前端user1拼接好数据对象Message
msg={id:1,userid:2,dstid:3,cmd:10,media:1,content:txt}
转化成json字符串jsonstr
jsonstr = JSON.stringify(msg)
通过websocket.send(jsonstr)发送
后端S在recvproc中接收收数据data
并做相应的逻辑处理dispatch(data)-转发给user2
user2通过websocket.onmessage收到消息后做解析并显示


###5.4.1 前端处理核心方法
前端所有的操作都在拼接数据
如何拼接?
```javascript
sendtxtmsg:function(txt){
//{id:1,userid:2,dstid:3,cmd:10,media:1,content:txt}
var msg =this.createmsgcontext();
//msg={"dstid":dstid,"cmd":cmd,"userid":userId()}
//选择某个好友/群的时候对dstid,cmd进行赋值
//userId()返回用户自己的id ,
// 从/chat/index.shtml?id=xx&token=yy中获得
//1文本类型
msg.media=1;msg.content=txt;
this.showmsg(userInfo(),msg);//显示自己发的文字
this.webSocket.send(JSON.stringify(msg))//发送
}

sendpicmsg:function(picurl){
    //{id:1,userid:2,dstid:3,cmd:10,media:4,
    // url:"http://www.baidu.com/a/log,jpg"}
    var msg =this.createmsgcontext();
    msg.media=4;
    msg.url=picurl;
    this.showmsg(userInfo(),msg)
    this.webSocket.send(JSON.stringify(msg))
}
sendaudiomsg:function(url,num){
    //{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
    var msg =this.createmsgcontext();
    msg.media=3;
    msg.url=url;
    msg.amount = num;
    this.showmsg(userInfo(),msg)
    console.log("sendaudiomsg",this.msglist);
    this.webSocket.send(JSON.stringify(msg))
}
```

##5.4.2 后端逻辑处理函数 func dispatch(data[]byte)
```cgo
func dispatch(data[]byte){
    //todo 解析data为message
    
    //todo根据message的cmd属性做相应的处理
    
}
func recvproc(node *Node) {
	for{
		_,data,err := node.Conn.ReadMessage()
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//todo 对data进一步处理
		dispatch(data)
		fmt.Printf("recv<=%s",data)
	}
}
```
###5.4.3 对端接收到消息后处理函数
```js
//初始化websocket的时候进行回调配置
this.webSocket.onmessage = function(evt){
     //{"data":"}",...}
     if(evt.data.indexOf("}")>-1){
         this.onmessage(JSON.parse(evt.data));
     }else{
         console.log("recv<=="+evt.data)
     }
 }.bind(this)
onmessage:function(data){
     this.loaduserinfo(data.userid,function(user){
         this.showmsg(user,data)
     }.bind(this))
 }

 //消息显示函数
showmsg:function(user,msg){
    var data={}
    data.ismine = userId()==msg.userid;
    //console.log(data.ismine,userId(),msg.userid)
    data.user = user;
    data.msg = msg;
    //vue 只需要修改数据结构即可完成页面渲染
    this.msglist = this.msglist.concat(data)
    //面板重置
    this.reset();
    var that =this;
    //滚动到新消息处
    that.timer = setTimeout(function(){
        window.scrollTo(0, document.getElementById("convo").offsetHeight);
        clearTimeout(that.timer)
    },100)
 }
```
###5.4.4 表情包简单逻辑
弹出一个窗口,
选择图片获得一个连接地址
调用sendpicmsg方法开始发送流程

##5.5 发送图片/拍照
弹出一个窗口,
选择图片,上传到服务器
获得一个链接地址
调用sendpicmsg方法开始发送流程
###5.5.1 界面处理技巧
```html
<input 
accept="image/gif,image/jpeg,,image/png" 
type="file" 
onchange="upload(this)" 
class='upload'/>
```
sendpicmsg方法开始发送流程
###5.5.2 upload前端实现
```javascript
function upload(dom){
        uploadfile("attach/upload",dom,function(res){
            if(res.code==0){//成功以后调用sendpicmsg
                vm.sendpicmsg(res.data)
            }
        })
    }
    
function uploadfile(uri,dom,callback){
    //H5新特性
    var formdata = new FormData();
    //获得一个文件dom.files[0]
    formdata.append("file",dom.files[0])
    //formdata.append("filetype",".png")//.mp3指定后缀
    
    var xhr = new XMLHttpRequest();//ajax初始化
    var url = "http://"+location.host+"/"+uri;
    //"http://127.0.0.1/attach/upload"
    xhr.open("POST",url, true);
    //成功时候回调
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && 
        xhr.status == 200) {
            //fn.call(this, JSON.parse(xhr.responseText));
            callback(JSON.parse(xhr.responseText))
        }
    };
    xhr.send(formdata);
}    
```
###5.5.2 upload后端实现
####存储到本地
```
func UploadLocal(writer http.ResponseWriter,
	request * http.Request){
	}
```
###存储到alioss
```
func UploadLocal(writer http.ResponseWriter,
	 request * http.Request){
}
如何安装 golang.org/x/time/rate
>cd $GOPATH/src/golang.org/x/
>git clone https://github.com/golang/time.git time

``` 
###5.6 发送语音
####5.6.1 采集语音
```javascript
navigator.mediaDevices.getUserMedia(
    {audio: true, video: true}
    ).then(successfunc).catch(errfunc);


navigator.mediaDevices.getUserMedia(
    {audio: true, video: false}
    ).then(function(stream)  {
              //请求成功
              this.recorder = new MediaRecorder(stream);
              this.recorder.start();
              this.recorder.ondataavailable = (event) => {
                  uploadblob("attach/upload",event.data,".mp3",res=>{
                      var duration = Math.ceil((new Date().getTime()-this.duration)/1000);
                      this.sendaudiomsg(res.data,duration);
                  })

                  stream.getTracks().forEach(function (track) {
                      track.stop();
                  });
                  this.showprocess = false
              }
              
          }.bind(this)).catch(function(err){
                mui.toast(err.msg)
                this.showprocess = false
            }.bind(this));
```
####5.6.2 上传语音
```javascript
function uploadblob(uri,blob,filetype,fn){
       var xhr = new XMLHttpRequest();
       xhr.open("POST","//"+location.host+"/"+uri, true);
       // 添加http头，发送信息至服务器时内容编码类型
       xhr.onreadystatechange = function() {
           if (xhr.readyState == 4 && (xhr.status == 200 || xhr.status == 304)) {
               fn.call(this, JSON.parse(xhr.responseText));
           }
       };
       var _data=[];
       var formdata = new FormData();
       formdata.append("filetype",filetype);
       formdata.append("file",blob)
       xhr.send(formdata);
   }
```

###5.7 实现群聊

####5.7.1 原理
分析群id,找到加了这个群的用户,把消息发送过去
方案一、
map<userid><qunid1,qunid2,qunid3>
优势是锁的频次低
劣势是要轮训全部map
```cgo
type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node,0)
```
方案二、
map<群id><userid1,userid2,userid3>
优势是找用户ID非常快
劣势是发送信息时需要根据userid获取node,锁的频次太高
```cgo
type Node struct {
	Conn *websocket.Conn
	//并行转串行,
	DataQueue chan []byte
}
//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node,0)
var comMap map[int64]set.Interface= make(map[int64]set.Interface,0)

```
####5.7.2 需要处理的问题
```javascript
1、当用户接入的时候初始化groupset
2、当用户加入群的时候刷新groupset
3、完成信息分发
```
###5.8 性能优化
1 锁的频次
2 json编码次数
3 静态资源分离
```cgo
###存储到alioss
```
func UploadOss(writer http.ResponseWriter,
	 request * http.Request){
}
如何安装
>go get github.com/aliyun/aliyun-oss-go-sdk/oss
 >golang.org/x/time/rate
>cd $GOPATH/src/golang.org/x/
>git clone https://github.com/golang/time.git time

``` 

```