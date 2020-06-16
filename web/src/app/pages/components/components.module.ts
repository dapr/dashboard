import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ComponentsComponent } from './components.component';

@NgModule({
  imports: [
    ThemeModule,
    RouterModule,
  ],
  declarations: [
    ComponentsComponent,
  ],
})
export class ComponentsModule { }
