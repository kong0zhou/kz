import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { DeployComponent } from './component/deploy/deploy.component';
import {ErrorStateMatcher,ShowOnDirtyErrorStateMatcher} from '@angular/material/core';
import { ReactiveFormsModule } from '@angular/forms';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations'
import { HttpClientModule } from '@angular/common/http';
import {
  MatTableModule, MatButtonModule, MatPaginatorModule,
  MatCardModule, MatDialogModule, MatSnackBarModule,
  MatProgressSpinnerModule, MatSelectModule, MatTabsModule,
  MatExpansionModule, MatFormFieldModule, MatInputModule,
  MatIconModule,MatRippleModule,MatDatepickerModule,MatNativeDateModule
} from '@angular/material';

@NgModule({
  declarations: [
    AppComponent,
    DeployComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
    MatButtonModule,
    HttpClientModule,
    // MatCheckboxModule,
    MatSelectModule,
    MatInputModule,
    MatFormFieldModule,
    // MatSnackBarModule,
    MatTableModule,
  ],
  providers: [
    // MatSnackBarModule,
    {provide: ErrorStateMatcher, useClass: ShowOnDirtyErrorStateMatcher},
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
