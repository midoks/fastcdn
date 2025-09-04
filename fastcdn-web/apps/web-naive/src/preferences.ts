import { defineOverridesPreferences } from '@vben/preferences';

/**
 * @description 项目配置文件
 * 只需要覆盖项目中的一部分配置，不需要的配置不用覆盖，会自动使用默认配置
 * !!! 更改配置后请清空缓存，否则可能不生效
 */
export const overridesPreferences = defineOverridesPreferences({
    // overrides
    app: {
        name: import.meta.env.VITE_APP_TITLE,
    },
    logo: {
        enable: true,
        fit: 'contain',
        source: 'https://unpkg.com/@vbenjs/static-source@0.1.7/source/logo-v1.webp',
        width: 180,
    },
    sidebar: {
        autoActivateChild: false,
        collapsed: false,
        collapsedButton: true,
        collapsedShowTitle: false,
        collapseWidth: 60,
        enable: true,
        expandOnHover: true,
        extraCollapse: false,
        extraCollapsedWidth: 60,
        fixedButton: true,
        hidden: false,
        mixedWidth: 160,
        width: 160,
    },
    theme: {
        builtinType: 'default',
        colorDestructive: 'hsl(348 100% 61%)',
        colorPrimary: 'hsl(212 100% 45%)',
        colorSuccess: 'hsl(144 57% 58%)',
        colorWarning: 'hsl(42 84% 61%)',
        mode: 'light',
        radius: '0.5',
        semiDarkHeader: false,
        semiDarkSidebar: false,
    },
});
