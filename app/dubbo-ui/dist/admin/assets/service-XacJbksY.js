import{d as w,k as S,r as d,O as C,u as T,b as R,f as m,c as a,t as e,e as n,h as r,P as p,K as D,o as l,F as V,x as M,v as _,I as v,y as f,z as h,H as O,J as P,m as B,_ as E}from"./index-tKheeObz.js";import{s as A}from"./service-hHA1p4Z2.js";import{S as Q,a as Y}from"./SearchUtil-6KhgO8Vi.js";/* empty css                                                               */import{g as j,a as q}from"./serverInfo-zEnFTwx2.js";import"./request-ovKoEMRQ.js";const z={class:"__container_app_service"},F={class:"statistic-icon-big"},H=w({__name:"service",setup(J){S(t=>({"402296cb":r(p)+"22","7369bb45":r(p)}));let s=d({info:{},report:{}}),g=d({info:{}});C(async()=>{o.tableStyle={scrollX:"100",scrollY:"calc(100vh - 600px)"};let t=(await j({})).data;g.info=(await q({})).data,s.info=t,s.report={providers:{icon:"carbon:branch",value:s.info.providers},consumers:{icon:"mdi:merge",value:s.info.consumers}}});const b=[{title:"idx",key:"idx"},{title:"服务",dataIndex:"serviceName",key:"serviceName",sorter:!0,width:"30%"},{title:"接口数",dataIndex:"interfaceNum",key:"interfaceNum",sorter:!0,width:"10%"},{title:"近 1min QPS",dataIndex:"avgQPS",key:"avgQPS",sorter:!0,width:"15%"},{title:"近 1min RT",dataIndex:"avgRT",key:"avgRT",sorter:!0,width:"15%"},{title:"近 1min 请求总量",dataIndex:"requestTotal",key:"requestTotal",sorter:!0,width:"15%"}],o=d(new Q([{label:"",param:"type",defaultValue:1,dict:[{label:"providers",value:1},{label:"consumers",value:2}],dictType:"BUTTON"},{label:"serviceName",param:"serviceName"}],A,b,{pageSize:4},!0));o.onSearch(),T();const x=R(),y=t=>{x.push("/resources/services/detail/"+t)};return D(B.SEARCH_DOMAIN,o),(t,K)=>{const I=n("a-statistic"),u=n("a-flex"),k=n("a-card"),N=n("a-button");return l(),m("div",z,[a(u,{wrap:"wrap",gap:"small",vertical:!1,justify:"start",align:"left"},{default:e(()=>[(l(!0),m(V,null,M(r(s).report,(i,c)=>(l(),_(k,{class:"statistic-card"},{default:e(()=>[a(u,{gap:"middle",vertical:!1,justify:"space-between",align:"center"},{default:e(()=>[a(I,{value:i.value,class:"statistic"},{prefix:e(()=>[a(r(v),{class:"statistic-icon",icon:"solar:target-line-duotone"})]),title:e(()=>[f(h(t.$t(c.toString())),1)]),_:2},1032,["value"]),O("div",F,[a(r(v),{icon:i.icon},null,8,["icon"])])]),_:2},1024)]),_:2},1024))),256))]),_:1}),a(Y,{"search-domain":o},{bodyCell:e(({column:i,text:c})=>[i.dataIndex==="serviceName"?(l(),_(N,{key:0,type:"link",onClick:L=>y(c)},{default:e(()=>[f(h(c),1)]),_:2},1032,["onClick"])):P("",!0)]),_:1},8,["search-domain"])])}}}),ee=E(H,[["__scopeId","data-v-89397748"]]);export{ee as default};