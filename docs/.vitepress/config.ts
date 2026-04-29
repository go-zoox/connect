import { defineConfig } from 'vitepress'

// 本地开发默认为 '/'。部署到 GitHub Pages 项目页时 CI 会设置 VITEPRESS_BASE=/connect/
const base = process.env.VITEPRESS_BASE ?? '/'

export default defineConfig({
  base,
  lang: 'zh-CN',
  title: 'GoZoox Connect',
  description:
    '轻量、强大的认证连接器（Auth Connect），支持 OAuth2、Doreamon、GitHub、飞书等。',
  cleanUrls: false,
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: '指南', link: '/guide/introduction', activeMatch: '^/guide/' },
      { text: '示例', link: '/examples/', activeMatch: '^/examples/' },
      { text: 'GitHub', link: 'https://github.com/go-zoox/connect' }
    ],
    sidebar: {
      '/guide/': [
        {
          text: '入门',
          items: [
            { text: '简介', link: '/guide/introduction' },
            { text: '安装', link: '/guide/installation' },
            { text: '快速开始', link: '/guide/quick-start' },
            { text: '部署形态', link: '/guide/architecture' }
          ]
        },
        {
          text: '使用手册',
          items: [
            { text: '命令行', link: '/guide/cli' },
            { text: '配置', link: '/guide/config' },
            { text: 'Docker 与编排', link: '/guide/docker' }
          ]
        }
      ],
      '/examples/': [
        {
          text: '示例',
          items: [
            { text: '概览', link: '/examples/' },
            { text: 'Doreamon 模式', link: '/examples/doreamon' },
            { text: 'GitHub OAuth', link: '/examples/github' },
            { text: '飞书 OAuth', link: '/examples/feishu' },
            { text: '上游代理（单站点）', link: '/examples/upstream' },
            { text: '无认证穿透', link: '/examples/none' }
          ]
        }
      ]
    },
    socialLinks: [{ icon: 'github', link: 'https://github.com/go-zoox/connect' }],
    footer: {
      message: '基于 MIT 许可证发布',
      copyright: 'Copyright © GoZoox'
    },
    search: { provider: 'local' },
    outline: { level: [2, 3] }
  }
})
