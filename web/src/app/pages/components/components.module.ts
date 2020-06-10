import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NbAccordionModule, NbCardModule, NbLayoutModule, NbIconModule, NbDialogModule } from '@nebular/theme';

import { ThemeModule } from '../../@theme/theme.module';
import { ComponentsComponent } from './components.component';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    NbLayoutModule,
    NbIconModule,
    NbDialogModule,
    NbAccordionModule,
    RouterModule
  ],
  declarations: [
    ComponentsComponent,
  ]
})
export class ComponentsModule { }
