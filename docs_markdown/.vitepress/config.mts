import {defineConfig} from 'vitepress'

export default defineConfig({
    title: 'Konfigo Docs',
    description: 'Documentation for the Konfigo project',
    outDir: '../docs', // Example: Output to /docs
    base: '/konfigo/', // Base path for deployment, adjust as needed 
    themeConfig: {
        nav: [
            {text: 'Home', link: '/'},
            {text: 'User Guide', link: '/guide/'},
            {text: 'Schema Guide', link: '/schema/'},
            {text: 'GitHub', link: 'https://github.com/ebogdum/konfigo'},
        ],
        sidebar: [
            {
                text: 'Getting Started',
                items: [
                    {text: 'Introduction', link: '/'},
                    {text: 'Installation', link: '/installation.md'},
                    {text: 'Quick Start', link: '/quick-start.md'}
                ]
            },
            {
                text: 'User Guide',
                items: [
                    {text: 'Overview', link: '/guide/'},
                    {text: 'CLI Reference', link: '/guide/cli-reference.md'},
                    {text: 'Environment Variables', link: '/guide/environment-variables.md'},
                    {text: 'Use Cases', link: '/guide/use-cases.md'}
                ]
            },
            {
                text: 'Schema Guide',
                items: [
                    {text: 'Overview', link: '/schema/'},
                    {text: 'Variables & Substitution', link: '/schema/variables.md'},
                    {text: 'Data Generation', link: '/schema/generation.md'},
                    {text: 'Data Transformation', link: '/schema/transformation.md'},
                    {text: 'Data Validation', link: '/schema/validation.md'},
                    {text: 'Advanced Concepts', link: '/schema/advanced.md'}
                ]
            }
        ]
    }
})