import { NgModule } from '@angular/core';
import { NbCardModule, NbLayoutModule, NbButtonModule, NbCheckboxModule, NbTabsetModule } from '@nebular/theme';

import { ThemeModule } from '../../@theme/theme.module';
import { DetailComponent } from '../detail/detail.component';


@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    NbLayoutModule,
    NbButtonModule,
    NbCheckboxModule,
    NbTabsetModule,
  ],
  declarations: [
    DetailComponent,
  ]
})
export class DetailModule { }