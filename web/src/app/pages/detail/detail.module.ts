import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
<<<<<<< HEAD
=======
import { NbCardModule, NbLayoutModule, NbButtonModule, NbCheckboxModule, NbTabsetModule } from '@nebular/theme';
import { ThemeModule } from '../../@theme/theme.module';
>>>>>>> develop
import { DetailComponent } from '../detail/detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';
import { LogsComponent } from './logs/logs.component';

@NgModule({
  imports: [
<<<<<<< HEAD
=======
    NbCardModule,
    ThemeModule,
    NbLayoutModule,
    NbButtonModule,
    NbCheckboxModule,
    NbTabsetModule,
>>>>>>> develop
    MonacoEditorModule,
    FormsModule,
  ],
  declarations: [
    DetailComponent,
<<<<<<< HEAD
    LogsComponent
=======
    LogsComponent,
>>>>>>> develop
  ],
})
export class DetailModule { }
