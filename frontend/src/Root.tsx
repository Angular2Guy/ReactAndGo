import * as React from 'react';
import { Router } from 'react-router';

import { browserHistory as history } from 'react-router';


// Routes
import configureRoutes from './configureRoutes';

const Root = ({history}) => (
    <Router history={history}>
      {configureRoutes()}
  </Router>
);

export default Root;

