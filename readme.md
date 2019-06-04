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
