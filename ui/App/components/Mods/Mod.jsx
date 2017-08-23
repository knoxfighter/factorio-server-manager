import React from 'react';

class Mod extends React.Component {
    constructor(props) {
        super(props);

        this.toggleUpdateStatus = this.toggleUpdateStatus.bind(this);
        this.removeVersionAvailableStatus = this.removeVersionAvailableStatus.bind(this);

        this.state = {
            newVersionAvailable: false,
            updateInProgress: false
        }
    }

    componentDidMount() {
        this.checkForNewVersion();
    }

    checkForNewVersion() {
        let this_class = this;
        //send AJAX that will check this
        $.ajax({
            url: "/api/mods/details",
            method: "POST",
            data: {
                mod_id: this.props.mod.name
            },
            dataType: "JSON",
            success: (data) => {
                let newest_release = JSON.parse(data.data).releases[0];
                console.log(newest_release);
                if(newest_release.version != this.props.mod.version) {
                    this_class.setState({
                        newVersionAvailable: true,
                        newVersion: {
                            downloadUrl: newest_release.download_url,
                            file_name: newest_release.file_name
                        }
                    });
                }
            },
            error: (jqXHR, status, err) => {
                console.log('api/mods/details', status, err.toString());
            }
        });
    }

    toggleUpdateStatus() {
        console.log("update Status changed");
        this.setState({
            updateInProgress: !this.state.updateInProgress
        })
    }

    removeVersionAvailableStatus() {
        this.setState({
            newVersionAvailable: false,
            newVersion: null
        })
    }

    render() {
        let modStatus;
        if(this.props.mod.enabled === false) {
            modStatus = <span className="label label-danger">Disabled</span>
        } else {
            modStatus = <span className="label label-success">Enabled</span>
        }

        let version;
        if(this.state.newVersionAvailable) {
            version = <span>{this.props.mod.version}
                <a className="btn btn-xs btn-default"
                   style={{
                       marginLeft: 10,
                       display: "inline-flex",
                       justifyContent: "center",
                       alignItems: "center",
                       width: 30,
                       height: 30
                   }}
                   href="#"
                   onClick={(event) => {
                       this.state.updateInProgress ? null : this.props.updateMod(event, this.toggleUpdateStatus, this.removeVersionAvailableStatus);
                   }}
                   data-download-url={this.state.newVersion.downloadUrl}
                   data-file-name={this.state.newVersion.file_name}
                >
                    {
                        this.state.updateInProgress ?
                            <div className='loader' style={{width: 15, height: 15, marginRight: 0, borderWidth: 3,}}></div>
                            :
                            <i className="fa fa-arrow-circle-up" title="Update Mod" style={{fontSize: "15pt"}}></i>
                    }
                </a>
            </span>;
        } else {
            version = this.props.mod.version;
        }

        return(
            <tr data-mod-name={this.props.mod.name}
                data-file-name={this.props.mod.file_name}
                data-mod-version={this.props.mod.version}
            >
                <td>{this.props.mod.title}</td>
                <td>{modStatus}</td>
                <td>{version}</td>
                <td>
                    <input className='btn btn-default btn-sm'
                        ref='modName'
                        type='submit'
                        value='Toggle'
                        onClick={this.props.toggleMod}
                    />

                    <input className="btn btn-danger btn-sm"
                        style={{marginLeft: 25}}
                        ref="modName"
                        type="submit"
                        value="Delete"
                        onClick={this.props.deleteMod}
                    />
                </td>
            </tr>
        )
    }
}

Mod.propTypes = {
    mod: React.PropTypes.object.isRequired,
    toggleMod: React.PropTypes.func.isRequired,
    deleteMod: React.PropTypes.func.isRequired,
    updateMod: React.PropTypes.func.isRequired,
};

export default Mod
