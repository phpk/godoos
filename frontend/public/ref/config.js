const getSystemConfig = (ifset = false) => {
  const configSetting = localStorage.getItem('GodoOS-config') || '{}';
  let config = JSON.parse(configSetting);

  config.version ??= '1.0.4';
  config.isFirstRun ??= false;
  config.lang ??= '';
  config.apiUrl ??= `${window.location.protocol}//${window.location.hostname}:56780`;
  config.userType ??= 'person';
  config.editorType ??= 'local';
  config.onlyoffice ??= { url: '', secret: '' };
  config.file ??= { isPwd: 0, pwd: '' };
  config.fileInputPwd ??= [];
  config.userInfo ??= {
    url: '',
    username: '',
    password: '',
    id: 0,
    nickname: '',
    avatar: '',
    email: '',
    phone: '',
    desc: '',
    job_number: '',
    work_place: '',
    hired_date: '',
    ding_id: '',
    role_id: 0,
    roleName: '',
    dept_id: 0,
    deptName: '',
    token: '',
    user_auths: '',
    user_shares: '',
    isPwd: false
  };

  config.isApp = window.go ? true : false;
  config.systemInfo ??= {};
  config.theme ??= 'light';
  config.storeType ??= localStorage.getItem('GodoOS-storeType') || 'local';
  config.storePath ??= '';
  config.netPath ??= '';
  config.netPort ??= '56780';
  config.background ??= {
    url: '/image/bg/bg6.jpg',
    type: 'image',
    color: 'rgba(30, 144, 255, 1)',
    imageList: [
      '/image/bg/bg1.jpg',
      '/image/bg/bg2.jpg',
      '/image/bg/bg3.jpg',
      '/image/bg/bg4.jpg',
      '/image/bg/bg5.jpg',
      '/image/bg/bg6.jpg',
      '/image/bg/bg7.jpg',
      '/image/bg/bg8.jpg',
      '/image/bg/bg9.jpg',
    ]
  };

  config.account ??= { ad: false, username: '', password: '' };
  if (config.userType === 'member') {
    config.account.ad = false;
  }
  config.storenet ??= { url: '', username: '', password: '', isCors: '' };
  config.webdavClient ??= { url: '', username: '', password: '' };
  config.dbInfo ??= { url: '', username: '', password: '', dbname: '' };
  config.chatConf ??= {
    checkTime: '15',
    first: '192',
    second: '168',
    thirdStart: '1',
    thirdEnd: '1',
    fourthStart: '2',
    fourthEnd: '254'
  };

  config.desktopList ??= [];
  config.menuList ??= [];
  config.token ??= generateRandomString(16);
  config = parseAiConfig(config);

  if (ifset) {
    setSystemConfig(config);
  }

  return config;
};

const setSystemConfig = (config) => {
  localStorage.setItem('GodoOS-config', JSON.stringify(config));
  localStorage.setItem('GodoOS-storeType', config.storeType);
};


function generateRandomString(length) {
  let result = '';
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const charactersLength = characters.length;
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}


const parseAiConfig = (config) => {
  if (!config.ollamaUrl) {
    config.ollamaUrl = `${window.location.protocol}//${window.location.hostname}:11434`
  }
  if (!config.ollamaDir) {
    config.ollamaDir = ''
  }
  if (!config.aiDir) {
    config.aiDir = ''
  }
  if (!config.aiUrl) {
    config.aiUrl = config.apiUrl
  }
  //openai
  if (!config.openaiUrl) {
    config.openaiUrl = 'https://api.openai.com/v1'
  }
  if (!config.openaiSecret) {
    config.openaiSecret = ""
  }
  //gitee
  if (!config.giteeSecret) {
    config.giteeSecret = ""
  }
  //cloudflare
  if (!config.cloudflareUserId) {
    config.cloudflareUserId = ""
  }
  if (!config.cloudflareSecret) {
    config.cloudflareSecret = ""
  }
  if (!config.deepseekSecret) {
    config.deepseekSecret = ""
  }
  if (!config.bigmodelSecret) {
    config.bigmodelSecret = ""
  }
  if (!config.volcesSecret) {
    config.volcesSecret = ""
  }
  if (!config.alibabaSecret) {
    config.alibabaSecret = ""
  }
  if (!config.groqSecret) {
    config.groqSecret = ""
  }
  if (!config.mistralSecret) {
    config.mistralSecret = ""
  }
  if (!config.anthropicSecret) {
    config.anthropicSecret = ""
  }
  if (!config.llamafamilySecret) {
    config.llamafamilySecret = ""
  }
  if (!config.siliconflowSecret) {
    config.siliconflowSecret = ""
  }
  return config
};