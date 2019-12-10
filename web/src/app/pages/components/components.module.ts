import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NbAccordionModule, NbCardModule, NbLayoutModule, NbIconModule, NbDialogModule } from '@nebular/theme';

import { ThemeModule } from '../../@theme/theme.module';
import { ComponentsComponent } from './components.component';
import { DataTableModule } from 'ng-angular8-datatable';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    DataTableModule,
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
