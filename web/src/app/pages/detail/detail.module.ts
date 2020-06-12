import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NbCardModule, NbLayoutModule, NbButtonModule, NbCheckboxModule, NbTabsetModule } from '@nebular/theme';
import { ThemeModule } from '../../@theme/theme.module';
import { DetailComponent } from '../detail/detail.component';
import { MonacoEditorModule } from 'ng-monaco-editor';

@NgModule({
  imports: [
    NbCardModule,
    ThemeModule,
    NbLayoutModule,
    NbButtonModule,
    NbCheckboxModule,
    NbTabsetModule,
    MonacoEditorModule,
    FormsModule,
  ],
  declarations: [
    DetailComponent,
  ],
})
export class DetailModule { }
