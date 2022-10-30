import { TestBed } from '@angular/core/testing';

import { InformationAPIService } from './information-api.service';

describe('InformationAPIService', () => {
  let service: InformationAPIService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(InformationAPIService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
