import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { MonacoEditorModule, MonacoProviderService } from 'ng-monaco-editor';

@NgModule({
  declarations: [AppComponent],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    AppRoutingModule,
    MonacoEditorModule.forRoot({
      baseUrl: 'lib',
      defaultOptions: {},
    }),
  ],
  providers: [
    {
      provide: MonacoProviderService,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
}
