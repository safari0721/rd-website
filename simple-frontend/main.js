(function(){
  const qs = new URLSearchParams(location.search)
  const apiBase = qs.get('api') || 'http://localhost:8080/api'
  document.getElementById('apiBase').textContent = apiBase

  const loginTab = document.getElementById('loginTab')
  const signupTab = document.getElementById('signupTab')
  const loginSec = document.getElementById('login')
  const signupSec = document.getElementById('signup')

  loginTab.addEventListener('click', ()=>{
    loginTab.classList.add('active'); signupTab.classList.remove('active')
    loginSec.classList.remove('hidden'); signupSec.classList.add('hidden')
  })
  signupTab.addEventListener('click', ()=>{
    signupTab.classList.add('active'); loginTab.classList.remove('active')
    signupSec.classList.remove('hidden'); loginSec.classList.add('hidden')
  })

  function showMsg(id, msg){ document.getElementById(id).textContent = msg }
  function showResult(obj){ document.getElementById('loginResult').textContent = JSON.stringify(obj, null, 2) }

  document.getElementById('signupBtn').addEventListener('click', async ()=>{
    const name = (document.getElementById('signup-name')).value.trim()
    const agent = (document.getElementById('signup-agent')).value.trim()
    const pass = (document.getElementById('signup-pass')).value
    showMsg('signupMsg','')
    if(name.length<1||name.length>100||agent.length<3||agent.length>64||pass.length<8||pass.length>128){
      showMsg('signupMsg','invalid input')
      return
    }
    try{
      const r = await fetch(apiBase + '/signup',{
        method:'POST', headers:{'Content-Type':'application/json'},
        body: JSON.stringify({ name:name, agent_id:agent, password:pass })
      })
      if(!r.ok){
        const e = await r.json().catch(()=>({error:'signup failed'}))
        showMsg('signupMsg', e.error || 'signup failed')
        return
      }
      showMsg('signupMsg','Signup successful')
    }catch(e){ showMsg('signupMsg','network error') }
  })

  document.getElementById('loginBtn').addEventListener('click', async ()=>{
    const agent = (document.getElementById('login-agent')).value.trim()
    const pass = (document.getElementById('login-pass')).value
    showMsg('loginMsg','')
    showResult('')
    if(agent.length<3||agent.length>64||pass.length<8||pass.length>128){
      showMsg('loginMsg','invalid input')
      return
    }
    try{
      const r = await fetch(apiBase + '/login',{
        method:'POST', headers:{'Content-Type':'application/json'},
        body: JSON.stringify({ agent_id:agent, password:pass })
      })
      if(!r.ok){
        const e = await r.json().catch(()=>({error:'login failed'}))
        showMsg('loginMsg', e.error || 'login failed')
        return
      }
      const data = await r.json()
      showMsg('loginMsg','Login successful')
      showResult(data)
    }catch(e){ showMsg('loginMsg','network error') }
  })
})()