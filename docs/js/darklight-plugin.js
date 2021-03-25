// Disclaimer: This file was extracted from https://github.com/boopathikumar018/docsify-darklight-theme, 
// and modified to fix an annoying bug that was breaking other plugins. A fix was already posted to GitHub
// issues (https://github.com/boopathikumar018/docsify-darklight-theme/issues/14), and as soon as it is 
// applied by the maintainer we can go back to fetching the plugin from a CDN, as usual.

const plugin = (hook, vm) => {
  
    var defaultConfig = {
      siteFont : "PT Sans",
      defaultTheme : 'dark',
      codeFontFamily : 'Roboto Mono, Monaco, courier, monospace',
      bodyFontSize : '17px',
      dark: {
        accent: '#42b983',
        toogleBackground : '#ffffff',
        background: '#091a28',
        textColor: '#b4b4b4',
        codeTextColor : '#ffffff',
        codeBackgroudColor : '#0e2233',
        borderColor : '#0d2538',
        blockQuoteColour : '#858585',
        highlightColor : '#d22778',
        sidebarSublink : '#b4b4b4',
        codeTypeColor : '#ffffff',
        coverBackground : 'linear-gradient(to left bottom, hsl(118, 100%, 85%) 0%,hsl(181, 100%, 85%) 100%)',
        toogleImage : 'url(https://cdn.jsdelivr.net/npm/docsify-darklight-theme@latest/icons/sun.svg)'
      },
      light: {
        accent: '#42b983',
        toogleBackground : '#091a28',
        background: '#ffffff',
        textColor: '#34495e',
        codeTextColor : '#525252',
        codeBackgroudColor : '#f8f8f8',
        borderColor : 'rgba(0, 0, 0, 0.07)',
        blockQuoteColor : '#858585',
        highlightColor : '#d22778',
        sidebarSublink : '#505d6b',
        codeTypeColor : '#091a28',
        coverBackground : 'linear-gradient(to left bottom, hsl(118, 100%, 85%) 0%,hsl(181, 100%, 85%) 100%)',
        toogleImage : 'url(https://cdn.jsdelivr.net/npm/docsify-darklight-theme@latest/icons/moon.svg)'
      }
    }
  
    let themeConfig = defaultConfig;
  
    if(vm.config.hasOwnProperty("darklightTheme")) {
      for (var [key, value] of Object.entries(vm.config.darklightTheme)) {
        if(key !== 'light' && key !== 'dark' && key !== 'defaultTheme') {
          themeConfig[key] = value;
        }
      }
  
      for (var [key, value] of Object.entries(themeConfig)) {
        if(key !== 'light' && key !== 'dark') {
          themeConfig[key] = value;
          document.documentElement.style.setProperty('--'+key , value);
        }
      }
  
      if(vm.config.darklightTheme.hasOwnProperty("dark")) {
        for (var [key, value] of Object.entries(vm.config.darklightTheme.dark)) {
          themeConfig.dark[key] = value
        }
      }
  
      if(vm.config.darklightTheme.hasOwnProperty("light")) {
        for (var [key, value] of Object.entries(vm.config.darklightTheme.light))
        themeConfig.light[key] = value 
      }
  
    } else {
      for (var [key, value] of Object.entries(themeConfig)) {
        if(key !== 'light' && key !== 'dark') {
          themeConfig[key] = value;
          document.documentElement.style.setProperty('--'+key , value);
        }
      }
    }
  
    if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
      themeConfig.defaultTheme = 'dark';
    } else if (window.matchMedia("(prefers-color-scheme: light)").matches) {
      themeConfig.defaultTheme = 'light';
    }
  
    var setTheme = (theme) => {
  
      localStorage.setItem('DARK_LIGHT_THEME', theme);
      themeConfig.defaultTheme = theme;
  
      if(theme == "light") {
        for (var [key, value] of Object.entries(themeConfig.light))
          document.documentElement.style.setProperty('--'+key , value)
      } else if ( theme == "dark") {
        for (var [key, value] of Object.entries(themeConfig.dark))
          document.documentElement.style.setProperty('--'+key , value)
      }
  
    }
    
    hook.afterEach(function(html, next) {
      var darkEl = `<div id="docsify-darklight-theme"><p>.</p></div>`
      html = `${darkEl}${html}`
      next(html)
    })
    
    hook.doneEach(function() {
      let savedTheme = localStorage.getItem('DARK_LIGHT_THEME')
      if ( savedTheme == "light" || savedTheme == "dark") {
        themeConfig.defaultTheme = savedTheme;
        setTheme(themeConfig.defaultTheme)
      } else {
        setTheme(themeConfig.defaultTheme);
      }
  
      const el = document.getElementById('docsify-darklight-theme')

      if (el !== null) {
        el.addEventListener('click', function() { themeConfig.defaultTheme === 'light' ? setTheme('dark') : setTheme('light')})
      }
        
    })
  }
    
  window.$docsify.plugins = [].concat(plugin, window.$docsify.plugins)