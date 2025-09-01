import { defineStore } from 'pinia';

// 预定义的主题颜色
export const THEMES = {
  blue: {
    name: '科技蓝',
    primary: '#409eff',
    primaryLight: '#79bbff',
    primaryDark: '#337ecc',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    background: '#f5f7fa',
    backgroundDark: '#1a1a1a'
  },
  green: {
    name: '自然绿',
    primary: '#67c23a',
    primaryLight: '#95d475',
    primaryDark: '#529b2e',
    gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    background: '#f0f9ff',
    backgroundDark: '#1a2332'
  },
  purple: {
    name: '优雅紫',
    primary: '#9c27b0',
    primaryLight: '#ba68c8',
    primaryDark: '#7b1fa2',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    background: '#faf5ff',
    backgroundDark: '#2d1b69'
  },
  orange: {
    name: '活力橙',
    primary: '#ff9800',
    primaryLight: '#ffcc02',
    primaryDark: '#f57c00',
    gradient: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)',
    background: '#fff8f0',
    backgroundDark: '#2d1b0f'
  },
  pink: {
    name: '樱花粉',
    primary: '#e91e63',
    primaryLight: '#f06292',
    primaryDark: '#c2185b',
    gradient: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
    background: '#fdf2f8',
    backgroundDark: '#2d1b20'
  }
};

export const useConfigStore = defineStore('config', {
    state: () => {
        return {
            mode: localStorage.getItem('theme-mode') || 'light',
            theme: localStorage.getItem('theme-color') || 'blue'
        };
    },
    getters: {
        currentTheme: (state) => THEMES[state.theme] || THEMES.blue,
        isDark: (state) => state.mode === 'dark',
        themeVars: (state) => {
            const theme = THEMES[state.theme] || THEMES.blue;
            return {
                '--primary-color': theme.primary,
                '--primary-light': theme.primaryLight,
                '--primary-dark': theme.primaryDark,
                '--gradient-bg': theme.gradient,
                '--bg-color': state.mode === 'dark' ? theme.backgroundDark : theme.background
            };
        }
    },
    actions: {
        toggle() {
            this.mode = this.mode === "dark" ? "light" : "dark";
            localStorage.setItem('theme-mode', this.mode);
            this.applyTheme();
        },
        setTheme(themeName) {
            if (THEMES[themeName]) {
                this.theme = themeName;
                localStorage.setItem('theme-color', themeName);
                this.applyTheme();
            }
        },
        applyTheme() {
            const root = document.documentElement;
            const theme = this.currentTheme;
            
            // 应用CSS变量
            Object.entries(this.themeVars).forEach(([key, value]) => {
                root.style.setProperty(key, value);
            });
            
            // 应用暗色模式类
            if (this.isDark) {
                document.body.classList.add('dark');
            } else {
                document.body.classList.remove('dark');
            }
        },
        initTheme() {
            this.applyTheme();
        }
    }
});