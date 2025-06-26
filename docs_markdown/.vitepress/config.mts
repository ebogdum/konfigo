import {defineConfig} from 'vitepress'

export default defineConfig({
    title: 'Konfigo Docs',
    description: 'Documentation for the Konfigo project',
    outDir: '../docs', // Example: Output to /docs
    base: '/konfigo/', // Base path for deployment, adjust as needed 
    themeConfig: {
        nav: [
            {text: 'Home', link: '/'},
            {text: 'Getting Started', link: '/getting-started/'},
            {text: 'User Guide', link: '/guide/'},
            {text: 'Schema Guide', link: '/schema/'},
            {text: 'Reference', link: '/reference/'},
            {text: 'GitHub', link: 'https://github.com/ebogdum/konfigo'},
        ],
        sidebar: [
            {
                text: 'Getting Started',
                items: [
                    {text: 'Introduction', link: '/'},
                    {text: 'Installation', link: '/getting-started/installation'},
                    {text: 'Quick Start', link: '/getting-started/quick-start'},
                    {text: 'Basic Concepts', link: '/getting-started/concepts'}
                ]
            },
            {
                text: 'User Guide',
                items: [
                    {text: 'Common Tasks', link: '/guide/'},
                    {text: 'Converting Formats', link: '/guide/format-conversion'},
                    {text: 'Merging Configurations', link: '/guide/merging'},
                    {text: 'Environment Variables', link: '/guide/environment-variables'},
                    {text: 'CLI Reference', link: '/guide/cli-reference'},
                    {text: 'Recipes & Examples', link: '/guide/recipes'}
                ]
            },
            {
                text: 'Schema Guide',
                items: [
                    {text: 'Schema Basics', link: '/schema/'},
                    {text: 'Variables & Substitution', link: '/schema/variables'},
                    {text: 'Validation', link: '/schema/validation'},
                    {text: 'Transformation', link: '/schema/transformation'},
                    {text: 'Data Generation', link: '/schema/generation'},
                    {text: 'Advanced Features', link: '/schema/advanced'}
                ]
            },
            {
                text: 'Reference',
                items: [
                    {text: 'Command Flags', link: '/reference/flags'},
                    {text: 'Configuration Options', link: '/reference/config'},
                    {text: 'Error Messages', link: '/reference/errors'},
                    {text: 'Best Practices', link: '/reference/best-practices'},
                    {text: 'Troubleshooting', link: '/reference/troubleshooting'},
                    {text: 'FAQ', link: '/reference/faq'}
                ]
            }
        ]
    }
})