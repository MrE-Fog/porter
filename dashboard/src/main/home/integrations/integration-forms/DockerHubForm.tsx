import React, { Component } from 'react';
import styled from 'styled-components';

import { Context } from '../../../../shared/Context';
import api from '../../../../shared/api';

import InputRow from '../../../../components/values-form/InputRow';
import SaveButton from '../../../../components/SaveButton';

type PropsType = {
};

type StateType = {
  registryURL: string,
  dockerEmail: string,
  dockerUsername: string,
  dockerPassword: string
};

export default class DockerHubForm extends Component<PropsType, StateType> {
  state = {
    registryURL: '',
    dockerEmail: '',
    dockerUsername: '',
    dockerPassword: ''
  }

  render() {
    return ( 
      <StyledDockerHubForm>
        <CredentialWrapper>
          <InputRow
            type='text'
            value={this.state.registryURL}
            setValue={(x: string) => this.setState({ registryURL: x })}
            label='📦 Registry URL'
            placeholder='ex: index.docker.io'
            width='100%'
          />
          <InputRow
            type='text'
            value={this.state.dockerEmail}
            setValue={(x: string) => this.setState({ dockerEmail: x })}
            label='✉️ Docker Email'
            placeholder='ex: captain@ahab.com'
            width='100%'
          />
          <InputRow
            type='text'
            value={this.state.dockerUsername}
            setValue={(x: string) => this.setState({ dockerUsername: x })}
            label='👤 Docker Username'
            placeholder='ex: whale_watcher_2000'
            width='100%'
          />
          <InputRow
            type='password'
            value={this.state.dockerPassword}
            setValue={(x: string) => this.setState({ dockerPassword: x })}
            label='🔒 Docker Password'
            placeholder='○ ○ ○ ○ ○ ○ ○ ○ ○'
            width='100%'
          />
        </CredentialWrapper>
        <SaveButton
          text='Save Changes'
          makeFlush={true}
          onClick={() => console.log('unimplemented')}
        />
      </StyledDockerHubForm>
    );
  }
}

const CredentialWrapper = styled.div`
  padding: 10px 40px 25px;
  background: #ffffff11;
  border-radius: 5px;
`;

const StyledDockerHubForm = styled.div`
  position: relative;
  padding-bottom: 75px;
`;