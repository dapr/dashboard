import { NgModule } from '@angular/core';
import { NbCardModule, NbLayoutModule, NbButtonModule } from '@nebular/theme';

import { ThemeModule } from '../../@theme/theme.module';
import { LogsComponent } from './logs.component';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    NbLayoutModule,
    NbButtonModule
  ],
  declarations: [
    LogsComponent,
  ]
})
export class LogsModule { }
