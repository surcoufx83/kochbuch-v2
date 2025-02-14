import { Component } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { SharedDataService } from './svc/shared-data.service';
import { ApiService } from './svc/api.service';
import { first } from 'rxjs';

@Component({
  selector: 'kb-root',
  templateUrl: './app.component.html',
  standalone: false,
  styleUrl: './app.component.scss'
})
export class AppComponent {

  constructor(
    apiService: ApiService,
    sharedDataService: SharedDataService,
    htmlTitleService: Title
  ) {
    sharedDataService.PageTitle.subscribe((t) => htmlTitleService.setTitle(t));
    apiService.get('').pipe(first()).subscribe((r) => {
      console.log(r)
    })
  }

}
