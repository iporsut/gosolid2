#!/usr/bin/env bash

curl -XPOST -v -H 'Content-Type: application/json' -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDUiLCJpc3MiOiJteWFwcCIsInN1YiI6InVzZXIiLCJhdWQiOlsibXlhcHBfdXNlcnMiXSwiZXhwIjoxNzU0MjA5OTA5LCJpYXQiOjE3NTQxMjM1MDl9.kmd4O9jlk6gWxk69GVk8Uomo3lXOxJNeDyGJ0ndbXXouI-Lc94SrLsxSScuES0tBJzk-EBYNZBhmlrHCihyf32SAX6RdpauSFjvCoqm8tS3Eva2CjHBCud16I53t4xNhg6XVd57A5offcHBDi6WiTOf3tT4YcsYrbWenzmWNeXrv4LSvoF4aHkWI2pJyHmm0mQ935DoDOxLVDw0_bEu796VTWkCx7hZkSWMYDxKcwbHMCOro1A6CY2AVXAUhvFwJR1fwPPenkfVMWYBJKimAPq1jpCgfa5hiOzUxVtSu6i7H6di_QT-HyogKdhMH2jca9iav8a7iAKCzFSZQPFQfDA" \
  localhost:8080/events -d \
  '{"name": "SOLID", "description": "SOLID concept in 1 hour", "number_of_tickets": 100, "start_date_time": "2025-08-03T10:00:00+07:00", "duration": 60}'
