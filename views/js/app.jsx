class Cores extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            cores: [],
            loading: true,
        };

        this.serverRequest = this.serverRequest.bind(this);
    }

    serverRequest() {
        $.get("/api/cores", res => {
            this.setState({
                cores: res,
                loading: false
            });
        });
    }

    componentDidMount() {
        this.serverRequest();
    }

    render() {

        let display = null;

        if (this.state.loading) {
            display = (
                <div class="ui active centered inline loader"></div>
            )
        } else {
            display = (
                <div className="ui items">
                    {this.state.cores.map(function (core, i) {
                        return <div className="item">
                            <div className="content">
                                <div className="header">
                                    {core.name}
                                </div>
                                <div className="meta">
                                    <span>{core.file}</span>
                                </div>
                                <div class="extra">
                                    <div class="ui right floated primary button">
                                        Download
                                        <i class="right chevron icon"></i>
                                    </div>
                                    <div class="ui label">MiSTer Core</div>
                                </div>
                            </div>
                        </div>
                    })}
                </div>
            )
        }
        return (
            <div className="ui container">
                <div className="ui stackable menu">
                    <div className="item">
                        MiSTer Bootstrap
                    </div>
                    <a className="item">Cores</a>
                </div>
                {display}
            </div>
        );
    }
}

class App extends React.Component {
    render() {
        return (<Cores />)
    }
}

ReactDOM.render(<App />, document.getElementById('app'));
