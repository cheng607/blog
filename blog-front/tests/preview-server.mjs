// Local-only fixture: never connects to a production backend.
import http from 'node:http'
import {readFile} from 'node:fs/promises'
import path from 'node:path'
const root=path.resolve('dist'), tags=[{id:1,name:'测试分类',count:12,description:'本地验证数据'}]
http.createServer(async(req,res)=>{
 const u=new URL(req.url,'http://127.0.0.1:4178')
 if(u.pathname.startsWith('/api/')){
  console.log(req.method,u.pathname,u.search)
  res.setHeader('Content-Type','application/json; charset=utf-8')
  const p=u.pathname.slice(4);let data
  if(p==='/site/summary')data={profile:{id:1,username:'本地测试博主',signature:'分页与草稿验证'},articleCount:12,tagCount:1,hollowCount:0}
  else if(p==='/tag')data=tags
  else if(p==='/tag/1')data=tags[0]
  else if(p==='/user')data=[{id:1,username:'本地测试博主',signature:'',avatar:''}]
  else if(p==='/hollow')data=[]
  else if(p==='/chat/status')data={configured:false}
  else if(p.startsWith('/auth/'))data={}
  else if(p==='/article/list'||p.startsWith('/article/tag/')){
   const page=Number(u.searchParams.get('page')||1),pageSize=5
   data={total:12,page,pageSize,items:Array.from({length:Math.max(0,Math.min(5,12-(page-1)*5))},(_,i)=>({id:12-(page-1)*5-i,title:`测试文章 ${12-(page-1)*5-i}`,excerpt:'服务端分页摘要，不包含全文。',tags,publishTime:'2026-11-09T08:00:00+08:00'}))}
  }else if(p==='/article'&&req.method==='POST'){res.statusCode=503;res.end(JSON.stringify({code:503,msg:'测试服务暂时不可用，草稿应保留'}));return}
  else if(/^\/article\/\d+$/.test(p))data={id:1,title:'测试文章',content:'# 阅读测试\n\n这是正文。\n\n```js\nconsole.log("hello")\n```',tags,publishTime:'2026-11-09T08:00:00+08:00'}
  else{res.statusCode=404;res.end(JSON.stringify({code:404,msg:'测试记录不存在'}));return}
  res.end(JSON.stringify({code:200,data}));return
 }
 try{
  let file=path.resolve(root,'.'+decodeURIComponent(u.pathname))
  if(!file.startsWith(root+path.sep)&&file!==root){res.statusCode=404;res.end();return}
  if(!path.extname(file))file=path.join(root,'index.html')
  let data=await readFile(file);const ext=path.extname(file)
  res.setHeader('Content-Type',({'.html':'text/html; charset=utf-8','.js':'text/javascript','.css':'text/css','.png':'image/png','.gif':'image/gif','.jpg':'image/jpeg'})[ext]||'application/octet-stream')
  if(ext==='.html')data=Buffer.from(data.toString().replace('<head>',`<head><script>sessionStorage.setItem('blog-admin-session','local-test-session')</script>`))
  res.end(data)
 }catch{res.statusCode=404;res.end('Not found')}
}).listen(4178,'127.0.0.1',()=>console.log('Local fixture http://127.0.0.1:4178'))
