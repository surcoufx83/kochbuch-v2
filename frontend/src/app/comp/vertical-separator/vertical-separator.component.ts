import { Component, Input } from '@angular/core';

@Component({
  selector: 'kb-vertical-separator',
  standalone: false,
  templateUrl: './vertical-separator.component.html',
  styleUrl: './vertical-separator.component.scss'
})
export class VerticalSeparatorComponent {

  @Input() height?: string = '4px';

}
