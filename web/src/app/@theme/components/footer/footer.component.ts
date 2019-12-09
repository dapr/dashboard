import { Component } from '@angular/core';

@Component({
  selector: 'ngx-footer',
  styleUrls: ['./footer.component.scss'],
  template: `
    <span class="created-by">Created with ♥ by <b><a href="http://dapr.io" target="_blank">Dapr</a></b> 2019</span>
  `,
})
export class FooterComponent {
}
