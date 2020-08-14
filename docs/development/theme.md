# Theming

## Angular Material Components
https://v9.material.angular.io/guide/theming

## Custom Components
https://v9.material.angular.io/guide/theming-your-components

## Theme

The themes, colors, and font are defined in `web/src/theme.scss`.

Additionally in this file, class definitions for the themes are available:
```scss
// theme.scss
...
.dashboard-light-theme {
    @include angular-material-theme($dashboard-light-theme);
    @include tables-theme($dashboard-light-theme);
    @include pages-theme($dashboard-light-theme);
    @include cards-theme($dashboard-light-theme);
    @include pages-component-theme($dashboard-light-theme);
}
...
```

`tables-theme`, `pages-theme`, etc. are custom components. For any custom component that needs to be styled according to a theme, a `<component>-theme.scss` file should be defined. This example styles `<a>` components as the primary theme color:
 
```scss
// tables-theme.scss
@import "~@angular/material/theming";

@mixin tables-theme($theme) {
    $primary: map-get($theme, primary);
    $accent: map-get($theme, accent);

    a {
        color: mat-color($primary);
        text-decoration: none;
    }
}
```