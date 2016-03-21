'use strict';

angular.module('aptWebApp')
  .controller('TableCtrl', function($scope, $http) {
    $scope.jt = {};
    $scope.jt.bsTableControl = {};

    $http.get('job/').then(response => {
      if (response) {
        $scope.jobs = response.data;

        $scope.jt.bsTableControl = {
          options: {
            data: $scope.jobs,
            rowStyle: function(row, index) {
              return {
                classes: 'none'
              };
            },
            cache: false,
            height: 400,
            striped: true,
            pagination: true,
            pageSize: 10,
            pageList: [5, 10, 25, 50, 100, 200],
            search: true,
            showColumns: true,
            showRefresh: true,
            minimumCountColumns: 2,
            clickToSelect: false,
            showToggle: false,
            maintainSelected: true,
            columns: [{
              field: 'state',
              checkbox: true
            }, {
              field: 'JobId',
              title: 'Job ID',
              align: 'center',
              valign: 'bottom',
              sortable: true
            }, {
              field: 'StartTime',
              title: 'Start Time',
              align: 'center',
              valign: 'middle',
              sortable: true
            }, {
              field: 'FrameKind',
              title: 'Framework',
              align: 'left',
              valign: 'top',
              sortable: true
            }, {
              field: 'FilterKind',
              title: 'Filter',
              align: 'left',
              valign: 'top',
              sortable: true
            }, {
              field: 'Status',
              title: 'Status',
              align: 'left',
              valign: 'top',
              sortable: true
            }]
          }
        };

      }
    });

  });
